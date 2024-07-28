package wikipedia

import (
	helpers "api/lib/helpers"
	"regexp"
	"strings"
)

func CleanArtistName(artistName string) string {
	regex := regexp.MustCompile(`(?i)(^The\s)?(?P<name>.*)\s*$`)
	matches := helpers.GetNamedMatchesForRegex(regex, artistName)

	if len(matches) != 0 {
		return strings.TrimSpace(matches["name"])
	}

	return strings.TrimSpace(artistName)
}

func CleanArtistNames(artistNames []string) []string {
	var cleanedArtistNames []string

	for _, artistName := range artistNames {
		cleanedArtistNames = append(cleanedArtistNames, CleanArtistName(artistName))
	}

	return cleanedArtistNames
}
