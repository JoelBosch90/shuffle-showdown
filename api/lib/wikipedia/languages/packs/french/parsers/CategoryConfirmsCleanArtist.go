package parsers

import (
	"api/lib/wikipedia/languages/packs/generic/helpers"
	"api/lib/wikipedia/models"
)

func CategoryConfirmsCleanedArtist(category models.Category, artistName string) bool {
	return CategoryConfirmsArtist(category, helpers.CleanArtistName(artistName))
}
