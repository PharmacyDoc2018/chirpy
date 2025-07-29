package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PharmacyDoc2018/chirpy/internal/auth"
	"github.com/google/uuid"
)

func handleWebhooks(mux *http.ServeMux, cfg *apiConfig) {
	mux.HandleFunc("POST /api/polka/webhooks", func(w http.ResponseWriter, req *http.Request) {
		token, err := auth.GetAPIKey(req.Header)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		if token != cfg.polkaSecret {
			fmt.Println("incorrect pokla secret in header")
			fmt.Println("cfg.polkaSecret:", cfg.polkaSecret)
			fmt.Println("token:", token)
			w.WriteHeader(401)
			return
		}
		webhook := polkaWebhook{}
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()
		err = decoder.Decode(&webhook)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		if webhook.Event != "user.upgraded" {
			w.WriteHeader(204)
			return
		}

		user_id, err := uuid.Parse(webhook.Data.UserId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}
		_, err = cfg.db.UpgradeToRedByID(req.Context(), user_id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(204)

	})
}
