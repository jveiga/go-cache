package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"unicode/utf8"

	"go-cache/pkg/cache"
)

// Api represents the api handler
type API struct {
	cache.Cacher
	host string
	port int
}

// New creates a new api
// in case this is extended with new functionalities,
// api configurations can be read here
func New(host string, port int, cacheExpiration time.Duration) *API {
	cache := cache.NewCache(cacheExpiration)
	return &API{Cacher: cache, host: host, port: port}
}

// Handle handles http requests to the cache api
func (a *API) Handle(w http.ResponseWriter, r *http.Request) {
	// remove / from path
	key := r.URL.Path[1:]
	// Check key size is not 0
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	// handle inserting in cache
	case "POST":
		// This should be changed to read max size in case it's exposed directly to internet
		// or downstream service should limit request or body size
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println("failed to read request body", err, string(body))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !isValidBody(body) {
			log.Println("received invalid request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		a.Insert(key, body)
		// Handle fetching from cache
	case "GET":
		value, found := a.Get(key)
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write(value)
	// case for UPDATE, HEAD, OPTIONS, PATCH
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func isValidBody(body []byte) bool {
	return len(body) > 0 && utf8.Valid(body)
}

// Run runs the web api
func (a *API) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.Handle)

	srv := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("%s:%d", a.host, a.port),
	}
	return srv.ListenAndServe()
}
