package english

import (
	helpers "api/lib/helpers"
	languageModels "api/lib/wikipedia/languages/models"
	formatters "api/lib/wikipedia/languages/packs/english/formatters"
)

var titleBasedFormatters = []languageModels.TitleBasedFormatter{formatters.FormatCleanTrackTitle, formatters.FormatTrackTitleWithArtistNamesSongSuffix, formatters.FormatTrackTitleWithOnlySongSuffix, formatters.FormatTrackTitleWithFirstrArtistNameSongSuffix}

func GetTitleSuggestions(trackTitle string, artistNames []string) []string {
	var titleSuggestions []string = []string{trackTitle}

	for _, formatter := range titleBasedFormatters {
		newSuggestion := formatter(trackTitle, artistNames)

		if newSuggestion != "" && !helpers.IncludesString(titleSuggestions, newSuggestion) {
			titleSuggestions = append(titleSuggestions, newSuggestion)
		}
	}

	return titleSuggestions
}
