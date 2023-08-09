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
	router.HandleFunc("/hello", handleRoot)

	http.Handle(
		"/metrics",
		middleware.New(registry, nil).
			WrapHandler("/metrics", promhttp.HandlerFor(
				registry,
				promhttp.HandlerOpts{}),
			))
	http.Handle("/", middleware.New(registry, nil).WrapHandler("/", router))
	err := http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server is closed!")
	} else if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	io.WriteString(w, "OK")
}
