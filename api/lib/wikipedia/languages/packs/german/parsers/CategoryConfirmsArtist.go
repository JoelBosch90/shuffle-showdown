package parsers

import (
	"api/lib/wikipedia/models"
	"regexp"
	"strings"
)

func CategoryConfirmsArtist(category models.Category, artistName string) bool {
	sanitizedArtistName := regexp.QuoteMeta(artistName)
	hyphenatedArtistName := strings.Replace(sanitizedArtistName, " ", "-", -1)
	regex := regexp.MustCompile("(?i)Kategorie:\\s*(Lied von " + sanitizedArtistName + "|" + hyphenatedArtistName + "-Lied|" + sanitizedArtistName + "\\s*(\\(Band\\))?)\\s*")
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
