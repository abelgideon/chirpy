package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	hits := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}
