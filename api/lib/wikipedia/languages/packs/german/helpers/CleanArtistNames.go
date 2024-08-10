package helpers

import (
	"api/lib/helpers"
	"strings"
)

func convertToLowerCaseExceptFirstCharaccter(artistName string) string {
	return strings.ToUpper(artistName[0:1]) + strings.ToLower(artistName[1:])
}

func CleanArtistName(artistName string) string {
	withoutSpace := strings.TrimSpace(artistName)
	withoutDiacritics := helpers.RemoveDiacritics(withoutSpace)
	withoutGuestArtists := convertToLowerCaseExceptFirstCharaccter(withoutDiacritics)

	return withoutGuestArtists
}

func CleanArtistNames(artistNames []string) []string {
	var cleanedArtistNames []string

	for _, artistName := range artistNames {
		cleanedArtistNames = append(cleanedArtistNames, CleanArtistName(artistName))
	}

	return cleanedArtistNames
}
