package wikipedia

import (
	wikipediaModels "api/lib/wikipedia/models"
	"regexp"
)

/**
 *	ParseLinkForSongSuffix returns the link title if the link has the (song) suffix.
 */
func ParseLinkForSongSuffix(link wikipediaModels.Link) string {
	regex := regexp.MustCompile(`(?i).*\(song\)`)
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}
