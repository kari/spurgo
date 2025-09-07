package main

import (
	"bufio"
	"errors"
	"math/rand/v2"
	"os"
	"strings"
)

// ErrNoMatch indicates that no matching lines were found
var ErrNoMatch = errors.New("no matching lines found")

// Sample reservoir samples a line from filename
// non-empty search limits sampling to lines that substring match
// https://gregable.com/2007/10/reservoir-sampling.html
func Sample(filename string, search string) (string, error) {
	search = strings.ToLower(search)

	data, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	var result string
	var line string
	i := 1

	for scanner.Scan() {
		line = scanner.Text()
		if search == "" || strings.Contains(strings.ToLower(line), search) {
			j := 1 + rand.IntN(i) // j = [1, i]
			if j <= 1 {
				result = line
			}
			i = i + 1
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if i == 1 {
		return "", ErrNoMatch
	}

	return result, nil
}
