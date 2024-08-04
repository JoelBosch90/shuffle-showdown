package formatters

import (
	"api/lib/wikipedia/models"
)

func FindLinkWithFirstArtistSuffix(page models.Page, trackTitle string, artistNames []string) string {
	return FindLinkWithArtistsSuffix(page, trackTitle, []string{artistNames[0]})
}
