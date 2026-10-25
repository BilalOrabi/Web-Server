package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func main() {

	ServeMux  := http.NewServeMux()
	path :=http.Dir(".")
	fileServer := http.FileServer(path)
	apiCfg := apiConfig{fileserverHits: atomic.Int32{},}
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer))

	ServeMux.Handle("assets/", fileServer)
	ServeMux.Handle("/app/", fsHandler)
	ServeMux.HandleFunc("GET /api/healthz", handlerReadiness)
	ServeMux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	ServeMux.HandleFunc("POST /api/reset", apiCfg.handlerReset)


	server := &http.Server{
		Handler: ServeMux,
		Addr: ":8080",
	}

	fmt.Println("Server starting on :8080...")
	server.ListenAndServe()
}

