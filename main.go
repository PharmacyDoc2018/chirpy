package main

import (
	"net/http"

	_ "github.com/lib/pq"
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
