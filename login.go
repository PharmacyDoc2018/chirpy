package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PharmacyDoc2018/chirpy/internal/auth"
	"github.com/PharmacyDoc2018/chirpy/internal/database"
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
			data, err := json.Marshal(returnErr{
				Error: "Incorrect email or password",
			})
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(401)
			w.Write(data)
			return
		}

		token, err := auth.MakeJWT(storedUser.ID, cfg.secret, maxTokenLifetime)
		if err != nil {
			data, err := json.Marshal(returnErr{
				Error: fmt.Sprint(err),
			})
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(400)
			w.Write(data)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		refreshTokenParams := database.CreateRefreshTokenParams{
			Token:     refreshToken,
			CreatedAt: time.Now().UTC(),
			UserID:    storedUser.ID,
			ExpiresAt: time.Now().UTC().Add(maxRefreshTokenLifetime),
		}

		_, err = cfg.db.CreateRefreshToken(req.Context(), refreshTokenParams)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		user := userResponse{
			ID:        storedUser.ID,
			CreatedAt: storedUser.CreatedAt,
			UpdatedAt: storedUser.UpdatedAt,
			Email:     storedUser.Email,
		}

		data, err := json.Marshal(loginResponse{
			Token:        token,
			RefreshToken: refreshToken,
			userResponse: user,
		})
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err)
			return
		}
		w.WriteHeader(200)
		w.Write(data)
	})
}
