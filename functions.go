package main

import (
	"database/sql"
	"fmt"
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

	cfg.platfrom = platform

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
	}
	cfg.db = database.New(db)

	return &cfg
}

func isCurse(word string) bool {
	curses := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	return slices.Contains(curses, strings.ToLower(word))
}

func filterProfanity(dirtyChirp string) string {
	censoredWord := "****"

	cleanChirpSplit := []string{}
	dirtyChirpSplit := strings.Split(dirtyChirp, " ")
	for _, word := range dirtyChirpSplit {
		if isCurse(word) {
			cleanChirpSplit = append(cleanChirpSplit, censoredWord)
		} else {
			cleanChirpSplit = append(cleanChirpSplit, word)
		}
	}
	cleanChirp := strings.Join(cleanChirpSplit, " ")

	return cleanChirp
}
