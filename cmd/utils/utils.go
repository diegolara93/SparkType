package utils

import (
	"github.com/tyler-smith/go-bip39/wordlists"
	"math/rand"
)

func GenerateWord(n int) string {
	words := wordlists.English

	s := ""
	for i := 0; i < n; i++ {
		s += words[rand.Intn(len(words))] + " "
		if i == n-1 {
			s += words[rand.Intn(len(words))]
		}
	}
	return s
}
