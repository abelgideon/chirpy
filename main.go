package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	apiConfig := apiConfig{}

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(fileHandler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("POST /api/reset", apiConfig.handlerReset)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()

}
