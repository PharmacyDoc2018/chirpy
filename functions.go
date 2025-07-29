package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/PharmacyDoc2018/chirpy/internal/database"
	"github.com/joho/godotenv"
)

func initapiConfig() *apiConfig {
	var cfg apiConfig

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polkaSecret := os.Getenv("POLKA_KEY")

	cfg.platfrom = platform
	cfg.secret = secret
	cfg.polkaSecret = polkaSecret

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
	}
	cfg.db = database.New(db)

	return &cfg
}

func decodeChirp(req *http.Request) (*chirp, *chirpError) {
	genericErrorWritten, err := json.Marshal(returnErr{Error: "something went wrong"})
	if err != nil {
		fmt.Println(err)
	}

	returnedChirp := chirp{}
	returnedDecodeError := chirpError{}

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err = decoder.Decode(&returnedChirp)
	if err != nil {
		returnedDecodeError.err = err
		returnedDecodeError.writtenErr = genericErrorWritten
		return &returnedChirp, &returnedDecodeError
	}

	return &returnedChirp, nil
}

func validateChirpLength(c *chirp) (bool, *chirpError) {
	const tooLongErrorMsg = "chirp is too long"
	tooLongErrorWritten, err := json.Marshal(returnErr{Error: tooLongErrorMsg})
	if err != nil {
		fmt.Println(err)
	}

	returnErr := chirpError{
		err:        errors.New(tooLongErrorMsg),
		writtenErr: tooLongErrorWritten,
	}

	if len(c.Body) > maxChirpLength {
		return false, &returnErr
	}

	return true, nil
}

func isCurse(word string) bool {
	curses := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	return slices.Contains(curses, strings.ToLower(word))
}

func filterProfanity(c *chirp) *chirp {
	const censoredWord = "****"

	cleanChirpSplit := []string{}
	dirtyChirpSplit := strings.Split(c.Body, " ")
	for _, word := range dirtyChirpSplit {
		if isCurse(word) {
			cleanChirpSplit = append(cleanChirpSplit, censoredWord)
		} else {
			cleanChirpSplit = append(cleanChirpSplit, word)
		}
	}
	cleanChirp := strings.Join(cleanChirpSplit, " ")
	c.Body = cleanChirp
	return c
}
