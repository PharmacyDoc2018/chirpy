package main

import (
	"slices"
	"strings"
)

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
