package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pliniogsnascimento/go-prometheus-sample/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	router := http.NewServeMux()
	router.HandleFunc("/", handleRoot)
	router.HandleFunc("/hello", handleHello)

	http.Handle(
		"/metrics",
		middleware.New(registry, nil).
			WrapHandler("/metrics", promhttp.HandlerFor(
				registry,
				promhttp.HandlerOpts{}),
			))

	// Add middlewares
	http.Handle("/", middleware.New(registry, nil).WrapHandler("/", router))
	http.Handle("/hello", middleware.New(registry, nil).WrapHandler("/hello", router))
	err := http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server is closed!")
	} else if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("Not found")
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "Not found")
}

func handleHello(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		log.Println("Request received on hello")
		io.WriteString(w, "Hello there!")
	default:
		log.Println("Request received on hello")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
