package utils

import "unicode"

func MixCase(text string) string {
	letters := []rune(text)
	for i, letter := range letters {
		if i%2 == 0 {
			letters[i] = unicode.ToUpper(letter)
		} else {
			letters[i] = unicode.ToLower(letter)
		}
	}
	return string(letters)
}
