package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

func main() {

	f, err := os.Create("../../../exp/examples/sort/unsorted.txt")
	if err != nil {
		log.Fatal(err)
	}

	words := []string{
		"Hello",
		"World",
		"Test",
		"Apple",
		"Banana",
		"Chair",
		"Table",
		"Dog",
		"Cat",
		"Run",
		"Jump",
		"Read",
		"Write",
		"Learn",
		"Play",
		"Friend",
		"Family",
		"Happy",
		"Sunny",
		"Clouds",
	}

	for range 1000000 {
		r := rand.Uint32()

		var sentence strings.Builder
		for i := range 10 {
			if i != 0 {
				sentence.WriteString(" ")
			}
			word := words[rand.Intn(len(words))]
			sentence.WriteString(word)

		}

		fmt.Fprintf(f, "%10d : %v\n", r, sentence.String())
	}

	f.Close()
}
