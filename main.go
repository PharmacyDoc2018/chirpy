package main

import (
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	//mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepathRoot+"/assets"))))
	mux.Handle("/assets/", http.FileServer(http.Dir(filepathRoot)))

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	srv.ListenAndServe()
}
