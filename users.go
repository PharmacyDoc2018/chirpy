package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PharmacyDoc2018/chirpy/internal/database"
	"github.com/google/uuid"
)

func handleResourceUsers(mux *http.ServeMux, cfg *apiConfig) {

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, req *http.Request) {
		genericErrorReturn := returnErr{
			Error: "Something went wrong",
		}

		w.Header().Set("Content-Type", "application/json")

		newUserReq := newUserRequest{}
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()
		err := decoder.Decode(&newUserReq)
		if err != nil {
			fmt.Printf("error decoding email: %s\n", err)
			w.WriteHeader(500)
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

		returnedNewUser := newUserResponse{
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

}
