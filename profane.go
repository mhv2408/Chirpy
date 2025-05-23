package main

import (
	"slices"
	"strings"
)

func removeProfaneWords(message string) string {
	profane_words := []string{"kerfuffle", "sharbert", "fornax"}
	word_slice := strings.Split(message, " ")
	for i := 0; i < len(word_slice); i++ {
		word := strings.ToLower(word_slice[i])
		if slices.Contains(profane_words, word) {
			word_slice[i] = "****"
		}
	}
	return strings.Join(word_slice, " ")

}
