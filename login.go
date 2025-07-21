package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PharmacyDoc2018/chirpy/internal/auth"
)

func handleLogin(mux *http.ServeMux, cfg *apiConfig) {
	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, req *http.Request) {
		loginInfo := loginRequest{}
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()
		err := decoder.Decode(&loginInfo)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			data, err := json.Marshal(returnErr{
				Error: fmt.Sprint(err),
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(data)
			return
		}

		storedUser, err := cfg.db.GetUserByEmail(req.Context(), loginInfo.Email)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			data, err := json.Marshal(returnErr{
				Error: fmt.Sprint(err),
			})
			if err != nil {
				fmt.Println(err)
			}
			w.Write(data)
			return
		}

		err = auth.CheckPasswordHash(loginInfo.Password, storedUser.HashedPassword)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			data, err := json.Marshal(returnErr{
				Error: "Incorrect email or password",
			})
			if err != nil {
				fmt.Println(err)
			}
			w.Write(data)
			return
		}

		data, err := json.Marshal(userResponse{
			ID:        storedUser.ID,
			CreatedAt: storedUser.CreatedAt,
			UpdatedAt: storedUser.UpdatedAt,
			Email:     storedUser.Email,
		})
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err)
		}
		w.WriteHeader(200)
		w.Write(data)
	})
}
