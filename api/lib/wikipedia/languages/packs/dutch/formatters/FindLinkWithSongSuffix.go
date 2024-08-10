package formatters

import (
	"api/lib/wikipedia/models"
	"regexp"
)

/**
 *	ParseLinkForSongSuffix returns the link title if the link has the (song) suffix.
 */
func parseLinkForSongSuffix(link models.Link) string {
	regex := regexp.MustCompile(`(?i).*\(lied\)`)
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}

func FindLinkWithSongSuffix(page models.Page, _ string, _ []string) string {
	var suggestedTitle string = ""

	for _, link := range page.Links {
		suggestedTitle = parseLinkForSongSuffix(link)
		if suggestedTitle != "" {
			return suggestedTitle
		}
	}

	return suggestedTitle
}
