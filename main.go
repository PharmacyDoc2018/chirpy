package main

import (
	"net/http"
)

func main() {
	cfg := initapiConfig()
	mux := initMux(cfg)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	srv.ListenAndServe()
}
