package main

import (
	_ "embed"
	"encoding/json"
	"math"
	"time"
)

//go:embed data/words_standard.json
var rawWords []byte

func words_standard(date time.Time) (word string, validWords map[string]struct{}) {
	var words struct {
		Solutions []string `json:"solutions"`
		Herrings  []string `json:"herrings"`
	}
	json.Unmarshal(rawWords, &words)

	validWords = make(map[string]struct{})
	for _, w := range words.Solutions {
		validWords[w] = struct{}{}
	}
	for _, w := range words.Herrings {
		validWords[w] = struct{}{}
	}

	t1 := time.Date(2021, time.Month(6), 19, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	days := t2.Sub(t1).Hours() / 24
	days = math.Abs(days)
	word = words.Solutions[int(days)%len(words.Solutions)]

	return
}
