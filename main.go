package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/atomicgo/cursor"
	"github.com/pterm/pterm"
)

func main() {
	var wordlen int
	var random bool
	var dateStr string
	flag.IntVar(&wordlen, "len", 5, "word length (only for random)")
	flag.BoolVar(&random, "random", false, "Random words (ignores website and uses wordlen)")
	flag.StringVar(&dateStr, "date", "", "Specify wordle website date (eg. 2021-09-15) (default is today)")
	flag.Parse()

	// Attempt to monospace BigText
	pterm.DefaultBigText.BigCharacters["i"] = `██████ 
  ██   
  ██   
  ██   
██████ `
	pterm.DefaultBigText.BigCharacters["I"] = pterm.DefaultBigText.BigCharacters["i"]

	var word string
	var validWords map[string]struct{}
	if random {
		word, validWords = words_random(wordlen)
		pterm.DefaultHeader.WithFullWidth().Println(fmt.Sprintf("Random Wordle\nLength: %d", wordlen))
	} else {
		dt := time.Now()
		if dateStr != "" {
			pt, terr := time.Parse("2006-01-02", dateStr)
			if terr != nil {
				log.Fatalln(terr)
			}
			dt = pt
		}
		word, validWords = words_standard(dt)
		pterm.DefaultHeader.WithFullWidth().Println("Standard Wordle\n" + dt.Format("Monday, January 2, 2006"))
	}
	pterm.DefaultCenter.Println("")

	styleWrong := pterm.NewStyle(pterm.FgGray)
	styleRight := pterm.NewStyle(pterm.FgGreen)
	styleHint := pterm.NewStyle(pterm.FgYellow)
	styleUnused := pterm.NewStyle(pterm.FgWhite, pterm.Bold)

	guessedLetters := make(map[string]*pterm.Style)
	alphabet := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	for _, alpha := range alphabet {
		guessedLetters[alpha] = styleUnused
	}

	guess := strings.Repeat("#", len(word))
	var attempts int
	var pletters []pterm.Letters
	exit := false
	giveup := false
	for !exit {
		cursor.ClearLinesUp(1)

		gl := strings.Split(guess, "")
		wl := strings.Split(word, "")

		for i := range wl {
			pstyle := styleWrong
			if strings.Contains(word, gl[i]) {
				pstyle = styleHint
				if gstyle := guessedLetters[gl[i]]; gstyle != styleRight {
					guessedLetters[gl[i]] = styleHint
				}
			}
			if gl[i] == wl[i] {
				pstyle = styleRight
				guessedLetters[gl[i]] = styleRight
			}
			if gstyle := guessedLetters[gl[i]]; gstyle == styleUnused {
				guessedLetters[gl[i]] = styleWrong
			}
			pletters = append(pletters, pterm.NewLettersFromStringWithStyle(gl[i], pstyle))
		}
		ptermLetters, _ := pterm.DefaultBigText.WithLetters(pletters...).Srender()
		pterm.DefaultCenter.Print(ptermLetters)

		var sb strings.Builder
		for _, alpha := range alphabet {
			sb.WriteString(guessedLetters[alpha].Sprint(strings.ToUpper(alpha)))
			sb.WriteString(" ")
		}
		pterm.DefaultCenter.Print(sb.String())

		if guess == word {
			exit = true
		} else {
			valid := false
			for !valid {
				fmt.Print("Guess (# to give up): ")
				fmt.Scanf("%s", &guess)
				if guess == "#" {
					exit = true
					giveup = true
					valid = true
				} else {
					_, valid = validWords[guess]
				}
				cursor.ClearLinesUp(1)
			}
		}

		attempts++
		pletters = nil
	}
	pterm.DefaultCenter.Print("")

	if giveup {
		attempts = 7
		pterm.DefaultCenter.Print("Try again later")
	}

	switch attempts {
	case 1:
		pterm.DefaultCenter.Print("Genius")
	case 2:
		pterm.DefaultCenter.Print("Magnificent")
	case 3:
		pterm.DefaultCenter.Print("Impressive")
	case 4:
		pterm.DefaultCenter.Print("Splendid")
	case 5:
		pterm.DefaultCenter.Print("Great")
	case 6:
		pterm.DefaultCenter.Print("Nice")
	default:
		pterm.DefaultCenter.Print("The word was " + word)
	}

	if giveup {
		os.Exit(1)
	}
}
