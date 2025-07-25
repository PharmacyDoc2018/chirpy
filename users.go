package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PharmacyDoc2018/chirpy/internal/auth"
	"github.com/PharmacyDoc2018/chirpy/internal/database"
	"github.com/google/uuid"
)

func handleResourceUsers(mux *http.ServeMux, cfg *apiConfig) {

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, req *http.Request) {
		genericErrorReturn := returnErr{
			Error: "Something went wrong",
		}

		w.Header().Set("Content-Type", "application/json")

		newUserReq := loginRequest{}
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()
		err := decoder.Decode(&newUserReq)
		if err != nil {
			fmt.Printf("error decoding login info: %s\n", err)
			w.WriteHeader(400)
			data, err := json.Marshal(genericErrorReturn)
			if err != nil {
				fmt.Println(err)
			}
			w.Write(data)
			return
		}

		params := database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			Email:     newUserReq.Email,
		}

		newUser, err := cfg.db.CreateUser(req.Context(), params)
		if err != nil {
			fmt.Println(err)
			errorReturn := returnErr{
				Error: fmt.Sprint(err),
			}
			data, err := json.Marshal(errorReturn)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(400)
			w.Write(data)
			return
		}

		hashedPassword, err := auth.HashPassword(newUserReq.Password)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		updatePassParams := database.UpdatePasswordParams{
			ID:             newUser.ID,
			HashedPassword: hashedPassword,
		}

		cfg.db.UpdatePassword(req.Context(), updatePassParams)

		returnedNewUser := userResponse{
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		}

		data, err := json.Marshal(returnedNewUser)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(201)
		w.Write(data)
	})

	mux.HandleFunc("PUT /api/users", func(w http.ResponseWriter, req *http.Request) {
		token, err := auth.GetBearerToken(req.Header)
		if err != nil {
			fmt.Println(err)
			data, err := json.Marshal(returnErr{
				Error: fmt.Sprint(err),
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

		userID, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil {
			fmt.Println(err)
			data, err := json.Marshal(returnErr{
				Error: fmt.Sprint(err),
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

		newLogin := loginRequest{}
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()
		err = decoder.Decode(&newLogin)
		if err != nil {
			fmt.Println(err)
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

		newEmail := newLogin.Email
		newHashedPassword, err := auth.HashPassword(newLogin.Password)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		params := database.UpdateEmailAndPasswordByIDParams{
			ID:             userID,
			Email:          newEmail,
			HashedPassword: newHashedPassword,
			UpdatedAt:      time.Now(),
		}

		storedUser, err := cfg.db.UpdateEmailAndPasswordByID(req.Context(), params)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		data, err := json.Marshal(userResponse{
			ID:        storedUser.ID,
			CreatedAt: storedUser.CreatedAt,
			UpdatedAt: storedUser.UpdatedAt,
			Email:     storedUser.Email,
		})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Write(data)

	})
}
