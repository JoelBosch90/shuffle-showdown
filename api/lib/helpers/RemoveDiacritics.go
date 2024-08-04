package helpers

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func RemoveDiacritics(text string) string {
	// Creates a transformer that:
	// 1. Normalizes the text to NFD form (e.g., transforms `é` to `e'`).
	// 2. Removes all non-spacing marks (e.g., transforms `e'` to `e`).
	// 3. Normalizes the text back to the NFC form, this time without diacritics.
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	withoutDiacritics, _, _ := transform.String(transformer, text)

	return withoutDiacritics
}
