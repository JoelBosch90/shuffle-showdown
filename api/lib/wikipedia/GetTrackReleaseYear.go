package wikipedia

import (
	helpers "api/lib/helpers"
	languages "api/lib/wikipedia/languages"
	languageModels "api/lib/wikipedia/languages/models"
	wikipediaModels "api/lib/wikipedia/models"
	"errors"
)

func getTitleSuggestions(languagePack languageModels.Pack, trackTitle string, artistNames []string) []string {
	var titleSuggestions []string = []string{trackTitle}

	for _, formatter := range languagePack.RedirectFormatters.TitleBasedFormatters {
		newSuggestion := formatter(trackTitle, artistNames)

		if newSuggestion != "" && !helpers.IncludesString(titleSuggestions, newSuggestion) {
			titleSuggestions = append(titleSuggestions, newSuggestion)
		}
	}

	return titleSuggestions
}

func getLinkSuggestions(languagePack languageModels.Pack, page wikipediaModels.Page, trackTitle string, artistNames []string) []string {
	var linkSuggestions []string = []string{trackTitle}

	for _, formatter := range languagePack.RedirectFormatters.LinkBasedFormatters {
		newSuggestion := formatter(page, trackTitle, artistNames)

		if newSuggestion != "" && !helpers.IncludesString(linkSuggestions, newSuggestion) {
			linkSuggestions = append(linkSuggestions, newSuggestion)
		}
	}

	return linkSuggestions
}

func trySuggestion(languagePack languageModels.Pack, trackTitle string, artistNames []string) (wikipediaModels.Response, uint) {
	if trackTitle == "" {
		return wikipediaModels.Response{}, 0
	}

	response, requestError := RequestTrackInfo(languagePack.Language, trackTitle)
	if requestError != nil {
		return response, 0
	}

	releaseYear, releaseYearError := GetConfirmedReleaseYear(languagePack, response, artistNames)
	if releaseYearError == nil {
		return response, releaseYear
	}

	return response, 0
}

func isPageMissing(response wikipediaModels.Response) bool {
	return len(response.Query.Pages) == 0 || response.Query.Pages[0].Missing
}

func GetTrackReleaseYear(trackTitle string, artistNames []string) (uint, error) {
	languageMaps := languages.Map

	for _, languagePack := range languageMaps {
		titleSuggestions := getTitleSuggestions(languagePack, trackTitle, artistNames)

		for _, suggestion := range titleSuggestions {
			response, releaseYear := trySuggestion(languagePack, suggestion, artistNames)

			if releaseYear != 0 {
				return releaseYear, nil
			}

			if isPageMissing(response) {
				continue
			}

			firstPage := response.Query.Pages[0]
			linkSuggestions := getLinkSuggestions(languagePack, firstPage, suggestion, artistNames)

			for _, suggestion := range linkSuggestions {
				response, releaseYear = trySuggestion(languagePack, suggestion, artistNames)

				if releaseYear != 0 {
					return releaseYear, nil
				}
			}
		}
	}

	return 0, errors.New("not a song page")
}
