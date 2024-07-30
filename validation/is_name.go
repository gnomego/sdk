package validation

import (
	"slices"
	"unicode"
)

func isNamePunctuation(r rune) bool {
	punc := []rune{'.', ',', '-', ' ', '\''}
	return slices.Contains(punc, r)
}

func IsNameValue(v string) bool {
	if len(v) == 0 {
		return true
	}

	for _, r := range v {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || isNamePunctuation(r) {
			continue
		}

		return false
	}

	return true
}
