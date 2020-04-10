package randomcode

import (
	"math/rand"
	"strings"
	"time"
)

var length = 3
var words = []string{}

// Generate returns a random code for reference
func Generate() string {
	parts := []string{}
	l := len(words)
	for i := 0; i < length; i++ {
		wordIdx := rand.Int() % l
		parts = append(parts, words[wordIdx])
	}

	return strings.Join(parts, "-")
}

func loadWord(word string) {
	words = append(words, word)
}

func init() {
	loadWords()
	rand.Seed(time.Now().UnixNano())
}
