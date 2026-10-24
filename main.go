package main

import (
	"fmt"
	"net/http"
)

func main() {

	ServeMux  := http.NewServeMux()
	path :=http.Dir(".")
	fileServer := http.FileServer(path)

	ServeMux.Handle("/assets/", fileServer)

	ServeMux.Handle("/app/", http.StripPrefix("/app", fileServer))
	
	ServeMux.HandleFunc("/healthz", func (w http.ResponseWriter, r *http.Request)  {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	})

	server := &http.Server{
		Handler: ServeMux,
		Addr: ":8080",
	}

	fmt.Println("Server starting on :8080...")
	server.ListenAndServe()
}

