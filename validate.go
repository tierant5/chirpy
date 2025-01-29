package main

import (
	"fmt"
	"strings"
)

func (c *Chirp) Validate() (string, bool) {
	const chirp_len = 140
	if len(c.Body) > chirp_len {
		msg := fmt.Sprintf("Chirp is longer than %v, got %v", chirp_len, len(c.Body))
		return msg, false
	}

	return replaceProfanity(c.Body), true
}

func replaceProfanity(msg string) string {
	const redactedWord = "****"
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(msg, " ")
	for i, word := range words {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			words[i] = redactedWord
		}
	}
	return strings.Join(words, " ")
}
