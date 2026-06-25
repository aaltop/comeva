package strings

import stdStrings "strings"

// SplitWords returns the passed strings split into words (separated by `sep`)
// and concatenated.
func SplitWords(sep string, sentences ...string) (words []string) {
	for _, sentence := range sentences {
		words = append(words, stdStrings.Split(sentence, sep)...)
	}
	return
}
