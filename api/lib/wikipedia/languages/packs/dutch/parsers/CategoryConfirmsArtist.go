package parsers

import (
	"api/lib/wikipedia/models"
	"regexp"
)

func CategoryConfirmsArtist(category models.Category, artistName string) bool {
	regex := regexp.MustCompile("(?i)Categorie:\\s*Nummer van\\s+" + regexp.QuoteMeta(artistName) + "\\s*")
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
