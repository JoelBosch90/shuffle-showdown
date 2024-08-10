package dutch

import (
	helpers "api/lib/helpers"
	languageModels "api/lib/wikipedia/languages/models"
	formatters "api/lib/wikipedia/languages/packs/dutch/formatters"
	models "api/lib/wikipedia/models"
)

var linkBasedFormatters = []languageModels.LinkBasedFormatter{formatters.FindLinkWithSongSuffix}

func GetLinkSuggestions(page models.Page, trackTitle string, artistNames []string) []string {
	var linkSuggestions []string = []string{trackTitle}

	for _, formatter := range linkBasedFormatters {
		newSuggestion := formatter(page, trackTitle, artistNames)

		if newSuggestion != "" && !helpers.IncludesString(linkSuggestions, newSuggestion) {
			linkSuggestions = append(linkSuggestions, newSuggestion)
		}
	}

	return linkSuggestions
}
