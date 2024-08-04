package wikipedia_languages_english

import (
	helpers "api/lib/helpers"
	"regexp"
	"strings"
)

func removeLeadingThe(artistName string) string {
	regex := regexp.MustCompile(`(?i)(^The\s)?(?P<name>.*)\s*$`)
	matches := helpers.GetNamedMatchesForRegex(regex, artistName)

	if len(matches) != 0 {
		return matches["name"]
	}

	return artistName
}

func CleanArtistName(artistName string) string {
	withoutSpace := strings.TrimSpace(artistName)
	withoutDiacritics := helpers.RemoveDiacritics(withoutSpace)
	withoutThe := removeLeadingThe(withoutDiacritics)

	return withoutThe
}

func CleanArtistNames(artistNames []string) []string {
	var cleanedArtistNames []string

	for _, artistName := range artistNames {
		cleanedArtistNames = append(cleanedArtistNames, CleanArtistName(artistName))
	}

	return cleanedArtistNames
}
