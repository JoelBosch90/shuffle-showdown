package wikipedia

import (
	wikipediaModels "api/lib/wikipedia/models"
	"regexp"
)

/**
 *	ParseLinkForCapitalization returns the link title if it matches the track title,
 * 	regardless of how the link is capitalized.
 */
func ParseLinkForCapitalization(link wikipediaModels.Link, trackTitle string) string {
	regex := regexp.MustCompile("(?i)^\\s*" + regexp.QuoteMeta(trackTitle) + "\\s*$")
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}
