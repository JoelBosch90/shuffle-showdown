package formatters

import (
	"api/lib/wikipedia/models"
)

func FindLinkWithAnyArtistSuffix(page models.Page, trackTitle string, artistNames []string) string {
	for _, artistName := range artistNames {
		suggestedTitle := FindLinkWithArtistsSuffix(page, trackTitle, []string{artistName})

		if suggestedTitle != "" {
			return suggestedTitle
		}
	}

	return ""
}
