package formatters

import (
	"api/lib/wikipedia/models"
	"regexp"
)

/**
 *	ParseLinkForCapitalization returns the link title if it matches the track title,
 * 	regardless of how the link is capitalized.
 */
func parseLinkForCapitalization(link models.Link, trackTitle string) string {
	sanitizedTrackTitle := regexp.QuoteMeta(trackTitle)
	withWildcardsForNonWordCharacters := regexp.MustCompile(`\w+`).ReplaceAllString(sanitizedTrackTitle, `\\w+`)
	regex := regexp.MustCompile("(?i)^\\s*" + withWildcardsForNonWordCharacters + "\\s*$")
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}

func FindLinkWithOtherCapitalization(page models.Page, trackTitle string, _ []string) string {
	var suggestedTitle string = ""

	for _, link := range page.Links {
		suggestedTitle = parseLinkForCapitalization(link, trackTitle)
		if suggestedTitle != "" && suggestedTitle != trackTitle {
			return suggestedTitle
		}
	}

	return suggestedTitle
}
