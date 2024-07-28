package wikipedia

import (
	wikipediaModels "api/lib/wikipedia/models"
	"regexp"
)

func ParseLinkForCapitalization(link wikipediaModels.Link, trackTitle string) string {
	regex := regexp.MustCompile("(?i)^\\s*" + regexp.QuoteMeta(trackTitle) + "\\s*$")
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}
