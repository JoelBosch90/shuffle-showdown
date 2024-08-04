package wikipedia

import (
	english "api/lib/wikipedia/languages/english"
	wikipediaModels "api/lib/wikipedia/models"
	"regexp"
	"strings"
)

func ParseLinkForArtistsSuffixWithCleaning(link wikipediaModels.Link, artistNames []string) string {
	return ParseLinkForArtistsSuffix(link, english.CleanArtistNames(artistNames))
}

/**
 *	ParseLinkForArtistsSuffix returns the link title if it matches the artist names,
 * 	regardless of how the link is capitalized.
 */
func ParseLinkForArtistsSuffix(link wikipediaModels.Link, artistNames []string) string {
	regex := regexp.MustCompile("(?i).*\\(" + regexp.QuoteMeta(strings.Join(artistNames, " and ")) + "\\s+song\\)")
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}
