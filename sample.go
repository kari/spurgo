package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

// sample reservoir samples a line from filename
// non-empty search limits sampling to lines that substring match
// https://gregable.com/2007/10/reservoir-sampling.html
func sample(filename string, search string) string {
	if filename == "" {
		filename = "data/vertauskuvat.txt"
	}
	rand.Seed(time.Now().UnixNano())
	file, _ := os.Open(filename)
	fscanner := bufio.NewScanner(file)
	i := 1
	var ret string

	for fscanner.Scan() {
		j := 1 + rand.Intn(i) // j = [1, i]
		if j <= 1 {
			ret = fscanner.Text()
		}
		i = i + 1
	}

	return ret
}
