package main

import "strings"

func isCurse(word string) bool {
	curses := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	for _, curse := range curses {
		if strings.ToLower(word) == curse {
			return true
		}
	}

	return false
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
