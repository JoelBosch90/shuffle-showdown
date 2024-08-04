package wikipedia_languages_english

import wikipediaModels "api/lib/wikipedia/models"

func CategoryConfirmsCleanedArtist(category wikipediaModels.Category, artistName string) bool {
	return CategoryConfirmsArtist(category, CleanArtistName(artistName))
}
