package main

import (
	"net/http"
)

func main() {
	mux := initMux()

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	srv.ListenAndServe()
}
