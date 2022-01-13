package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"time"
)

//go:embed data/words_alpha.txt.gz
var wordsCompressed []byte

func words_random(wordlen int) (word string, validWords map[string]struct{}) {
	var wordlist []string
	validWords = make(map[string]struct{})
	wordsin, _ := gzip.NewReader(bytes.NewReader(wordsCompressed))
	scanner := bufio.NewScanner(wordsin)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == wordlen {
			wordlist = append(wordlist, word)
			validWords[word] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading file:", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rIdx := r.Intn(len(wordlist))
	word = wordlist[rIdx]

	return
}
