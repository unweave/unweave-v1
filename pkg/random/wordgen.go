package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func getRandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func GenerateRandomAdjectiveNounTriplet() string {
	nl := len(nouns)
	al := getRandomInt(0, len(adjectives))
	triplet := fmt.Sprintf("%s-%s-%s", adjectives[al], nouns[getRandomInt(0, nl)], nouns[getRandomInt(0, nl)])
	return triplet
}

func GenerateRandomWord() string {
	return wordList[getRandomInt(0, len(wordList))]
}

func GenerateRandomPhrase(numWords int, separator string) string {
	results := make([]string, numWords)

	for i := 0; i < numWords; i++ {
		results[i] = GenerateRandomWord()
	}
	return strings.Join(results, separator)
}

func GenerateRandomEmoji() string {
	t := rand.Intn(3)
	switch t {
	case 0:
		return natureEmojis[getRandomInt(0, len(natureEmojis))]
	case 1:
		return foodEmojis[getRandomInt(0, len(foodEmojis))]
	default:
		return transportEmojis[getRandomInt(0, len(transportEmojis))]
	}
}
