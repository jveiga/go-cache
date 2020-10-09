package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"go-cache/pkg/cache"
)

// used for testing
var defaultHost = "localhost"
var defaultPort = 5000

// util function to create and start server
func setupServer() *http.Server {
	c := cache.NewCache(10 * time.Second)
	api := &API{Cacher: c}
	mux := http.NewServeMux()
	mux.HandleFunc("/", api.Handle)
	srv := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("%s:%d", defaultHost, defaultPort),
	}
	go log.Println(srv.ListenAndServe())
	time.Sleep(100 * time.Millisecond)
	return srv
}

func TestGetMissingKey(t *testing.T) {
	srv := setupServer()
	defer srv.Close()
	time.Sleep(time.Second)
	url := fmt.Sprintf("http://%s:%d/key", defaultHost, defaultPort)
	log.Println(url)

	resp, err := http.Get(url)

	if err != nil {
		t.Fatalf("request should not fail %s", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Incorrect status code, expecting 404, got %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("reading response body should not fail", err)
	}
	defer resp.Body.Close()
	if len(body) != 0 {
		t.Fatal("Expecting body to be empty")
	}
}

func TestGetInvalidKeyEmptyKey(t *testing.T) {
	srv := setupServer()
	defer srv.Close()
	key := ""
	url := fmt.Sprintf("http://%s:%d/%s", defaultHost, defaultPort, key)

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("request should not fail %s", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Incorrect status code, expecting 200, got %d", resp.StatusCode)
	}
}

func TestPostInvalidKeyEmptyKey(t *testing.T) {
	srv := setupServer()
	defer srv.Close()
	key, value := "", []byte("somevalue")
	url := fmt.Sprintf("http://%s:%d/%s", defaultHost, defaultPort, key)

	resp, err := http.Post(url, "", bytes.NewBuffer(value))
	if err != nil {
		t.Fatalf("request should not fail %s", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Incorrect status code, expecting 200, got %d", resp.StatusCode)
	}
}

func TestGetExistingKey(t *testing.T) {
	srv := setupServer()
	defer srv.Close()
	key, value := "somekey", []byte("somevalue")
	url := fmt.Sprintf("http://%s:%d/%s", defaultHost, defaultPort, key)

	// prime cache
	resp, err := http.Post(url, "", bytes.NewBuffer(value))
	if err != nil {
		t.Fatal("request should not fail")
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Incorrect status code, expecting 200, got %d", resp.StatusCode)
	}

	// get from cache
	getResp, err := http.Get(url)
	if err != nil {
		t.Fatal("request should not fail", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("Incorrect status code, expecting 200, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		log.Fatal("reading response body should not fail", err)
	}
	defer resp.Body.Close()
	if !bytes.Equal(value, body) {
		t.Fatalf("Expecting %s, got '%s'", string(value), string(body))
	}
}

func TestPostNonUTF8Body(t *testing.T) {
	srv := setupServer()
	defer srv.Close()
	key, value := "key", []byte("\xc5z")
	url := fmt.Sprintf("http://%s:%d/%s", defaultHost, defaultPort, key)

	resp, err := http.Post(url, "", bytes.NewBuffer(value))
	if err != nil {
		t.Fatalf("request should not fail %s", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Incorrect status code, expecting 200, got %d", resp.StatusCode)
	}
}

func TestPostEmptyBody(t *testing.T) {
	srv := setupServer()
	defer srv.Close()
	key, value := "key", []byte{}
	url := fmt.Sprintf("http://%s:%d/%s", defaultHost, defaultPort, key)

	resp, err := http.Post(url, "", bytes.NewBuffer(value))
	if err != nil {
		t.Fatalf("request should not fail %s", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Incorrect status code, expecting 200, got %d", resp.StatusCode)
	}
}
