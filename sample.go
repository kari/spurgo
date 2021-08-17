package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

// sample reservoir samples a line from filename
// non-empty search limits sampling to lines that substring match
// https://gregable.com/2007/10/reservoir-sampling.html
func sample(filename string, search string) string {
	if filename == "" {
		filename = "data/vertauskuvat.txt" // FIXME: embed using go:embed
	}
	search = strings.ToLower(search)

	rand.Seed(time.Now().UnixNano())
	file, _ := os.Open(filename)
	fscanner := bufio.NewScanner(file)
	i := 1
	var ret string
	var line string

	for fscanner.Scan() {
		line = fscanner.Text()
		if strings.Contains(strings.ToLower(line), search) {
			j := 1 + rand.Intn(i) // j = [1, i]
			if j <= 1 {
				ret = line
			}
			i = i + 1
		}
	}

	return ret
}
