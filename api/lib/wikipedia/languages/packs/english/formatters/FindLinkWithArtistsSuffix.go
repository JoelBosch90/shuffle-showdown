package formatters

import (
	"api/lib/wikipedia/languages/packs/english/helpers"
	"api/lib/wikipedia/models"
	"regexp"
	"strings"
)

/**
 *	ParseLinkForArtistsSuffix returns the link title if it matches the artist names,
 * 	regardless of how the link is capitalized.
 */
func parseLinkForArtistsSuffix(link models.Link, artistNames []string) string {
	regex := regexp.MustCompile("(?i).*\\(" + regexp.QuoteMeta(strings.Join(artistNames, " and ")) + "\\s+song\\)")
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}

func parseLinkForArtistsSuffixWithCleaning(link models.Link, artistNames []string) string {
	return parseLinkForArtistsSuffix(link, helpers.CleanArtistNames(artistNames))
}

func FindLinkWithArtistsSuffix(page models.Page, trackTitle string, artistNames []string) string {
	var suggestedTitle string = ""

	for _, link := range page.Links {
		suggestedTitle = parseLinkForArtistsSuffix(link, artistNames)
		if suggestedTitle != "" {
			return suggestedTitle
		}
		suggestedTitle = parseLinkForArtistsSuffixWithCleaning(link, artistNames)
		if suggestedTitle != "" {
			return suggestedTitle
		}
	}

	return suggestedTitle
}
