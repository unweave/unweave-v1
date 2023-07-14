//nolint:gosec
package random

import (
	"encoding/base64"
	"math/rand"
)

const (
	letterBytes          = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	safeLowerLetterBytes = "0123456789bcdfghjklmnpqrstvwxyz" // vowels removed so we can't get swears
	letterIdxBits        = 6                                 // 6 bits to represent a letter index
	letterIdxMask        = 1<<letterIdxBits - 1              // All 1-bits, as many as letterIdxBits
	letterIdxMax         = 63 / letterIdxBits                // # of letter indices fitting in 63 bits
)

func GenerateRandomString(n int) (string, error) {
	buffer := make([]byte, n)
	// A src.Int63() generator is used to give us 63 random bits, enough to cover the letterIdxMax characters
	for letterIdx, cache, remain := n-1, rand.Int63(), letterIdxMax; letterIdx >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			buffer[letterIdx] = letterBytes[idx]
			letterIdx--
		}

		cache >>= letterIdxBits
		remain--
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

// GenerateRandomLower creates a string of n length
// that is entirely lowercase letters and numbers.
// It's considered 'safe' because the alphabet has
// vowels removed so that we cannot accidentally
// generate swear words.
func GenerateRandomLower(n int) string {
	buffer := make([]byte, n)
	// A src.Int63() generator is used to give us 63 random bits, enough to cover the letterIdxMax characters
	for letterIdx, cache, remain := n-1, rand.Int63(), letterIdxMax; letterIdx >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(safeLowerLetterBytes) {
			buffer[letterIdx] = safeLowerLetterBytes[idx]
			letterIdx--
		}

		cache >>= letterIdxBits
		remain--
	}

	return string(buffer)
}
