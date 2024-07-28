package wikipedia

import (
	wikipediaModels "api/lib/wikipedia/models"
	"regexp"
	"strings"
)

func ParseLinkForArtistsSuffixWithCleaning(link wikipediaModels.Link, artistNames []string) string {
	return ParseLinkForArtistsSuffix(link, CleanArtistNames(artistNames))
}

func ParseLinkForArtistsSuffix(link wikipediaModels.Link, artistNames []string) string {
	regex := regexp.MustCompile("(?i).*\\(" + regexp.QuoteMeta(strings.Join(artistNames, " and ")) + "\\s+song\\)")
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}
