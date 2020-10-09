package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"time"

	"go-cache/pkg/api"
)

// flags for profiling web server
var devMode = flag.Bool("dev", false, "starts profiling web server")
var devModeHost = flag.String("profiler-host", "localhost:6000", "profiling web server host")

// flags for configuring web server
var host = flag.String("host", "127.0.0.1", "host:port on which the server runs")
var port = flag.Int("port", 8000, "port on which the caching webserver listens")
var cacheDuration = flag.Duration("duration", 30*time.Minute, "cache time duration")

func init() {
	flag.Parse()
	if *devMode {
		go func() {
			log.Println(http.ListenAndServe(*devModeHost, nil))
		}()
	}
}

func main() {
	host, port, err := parseFlags()
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}
	instance := api.New(host, port, *cacheDuration)
	if err := instance.Run(); err != nil {
		log.Println(err)
	}
}

func parseFlags() (string, int, error) {
	if *port > 65535 {
		return "", 0, errors.New("Invalid port value")
	}
	u, err := url.Parse(*host)
	if err != nil {
		return "", 0, fmt.Errorf("Invalid host given %w", err)
	}

	return u.Hostname(), *port, nil
}
