package helpers

import (
	"api/lib/helpers"
	generic "api/lib/wikipedia/languages/packs/generic/helpers"
)

func CleanArtistNames(artistNames []string) []string {
	var cleanedArtistNames []string

	for _, artistName := range artistNames {
		cleanedArtistNames = append(cleanedArtistNames, helpers.ConvertToLowerCaseExceptFirstCharacter(generic.CleanArtistName(artistName)))
	}

	return cleanedArtistNames
}
