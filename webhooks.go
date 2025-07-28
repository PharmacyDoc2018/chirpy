package main

import "net/http"

func handleWebhooks(mux *http.ServeMux, cfg *apiConfig) {
	mux.HandleFunc("POST /api/polka/webhooks", func(w http.ResponseWriter, req *http.Request) {
		//
	})
}
