package main

import "net/http"

const filepathRoot = "."
const port = "8080"

func initMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/assets/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})

	return mux
}
