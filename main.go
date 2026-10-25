package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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
	ServeMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	ServeMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	ServeMux.HandleFunc("POST /api/validate_chirp", handlerChripsValidate)

	server := &http.Server{
		Handler: ServeMux,
		Addr: ":8080",
	}

	fmt.Println("Server starting on :8080...")
	server.ListenAndServe()
}

