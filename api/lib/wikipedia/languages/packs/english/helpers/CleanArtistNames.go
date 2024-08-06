package helpers

import (
	"api/lib/helpers"
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

func removeLeadingAsterisk(artistName string) string {
	regex := regexp.MustCompile(`(?i)^\*\s*(?P<name>.*)\s*$`)
	matches := helpers.GetNamedMatchesForRegex(regex, artistName)

	if len(matches) != 0 {
		return matches["name"]
	}

	return artistName
}

func removeLeadingPrefix(artistName string) string {
	prefixes := []string{"Ms.", "Mr.", "Ms.", "Dr.", "Prof.", "Prof"}
	prefixesRegex := strings.Join(prefixes, "|")
	regex := regexp.MustCompile(`(?i)^(` + prefixesRegex + `)\s*(?P<name>.*)\s*$`)
	matches := helpers.GetNamedMatchesForRegex(regex, artistName)

	if len(matches) != 0 {
		return matches["name"]
	}

	return artistName
}

func removeTrailingGuestArtists(artistName string) string {
	regex := regexp.MustCompile(`(?i)(?P<name>.*)(\s+featuring\s+.*|\s+ft\..*|\s+&\s+.*|\s+with\s+.*|\s+and\s+.*)$`)
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
	withoutAsterisk := removeLeadingAsterisk(withoutThe)
	withoutPrefix := removeLeadingPrefix(withoutAsterisk)
	withoutGuestArtists := removeTrailingGuestArtists(withoutPrefix)

	return withoutGuestArtists
}

func CleanArtistNames(artistNames []string) []string {
	var cleanedArtistNames []string

	for _, artistName := range artistNames {
		cleanedArtistNames = append(cleanedArtistNames, CleanArtistName(artistName))
	}

	return cleanedArtistNames
}
