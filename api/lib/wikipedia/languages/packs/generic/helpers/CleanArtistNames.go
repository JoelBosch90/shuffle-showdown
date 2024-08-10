package helpers

import (
	"api/lib/helpers"
	"strings"
)

func CleanArtistName(artistName string) string {
	withoutSpace := strings.TrimSpace(artistName)
	withoutDiacritics := helpers.RemoveDiacritics(withoutSpace)
	withoutGuestArtists := helpers.ConvertToLowerCaseExceptFirstCharacter(withoutDiacritics)

	return withoutGuestArtists
}

func CleanArtistNames(artistNames []string) []string {
	var cleanedArtistNames []string

	for _, artistName := range artistNames {
		splitNames := strings.Split(artistName, "&")

		for _, splitName := range splitNames {
			cleanedArtistNames = append(cleanedArtistNames, CleanArtistName(splitName))
		}
	}

	return cleanedArtistNames
}
