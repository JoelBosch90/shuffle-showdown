package wikipedia

import (
	languages "api/lib/wikipedia/languages"
	languageModels "api/lib/wikipedia/languages/models"
	wikipediaModels "api/lib/wikipedia/models"
	"errors"
)

func trySuggestion(languagePack languageModels.Pack, trackTitle string, artistNames []string) (wikipediaModels.Response, uint) {
	if trackTitle == "" {
		return wikipediaModels.Response{}, 0
	}

	response, requestError := RequestTrackInfo(languagePack.Language, trackTitle)
	if requestError != nil {
		return response, 0
	}

	releaseYear, releaseYearError := languagePack.GetReleaseYear(response, artistNames)
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
		titleSuggestions := languagePack.GetTitleSuggestions(trackTitle, artistNames)

		for _, suggestion := range titleSuggestions {
			response, releaseYear := trySuggestion(languagePack, suggestion, artistNames)

			if releaseYear != 0 {
				return releaseYear, nil
			}

			if isPageMissing(response) {
				continue
			}

			firstPage := response.Query.Pages[0]
			linkSuggestions := languagePack.GetLinkSuggestions(firstPage, suggestion, artistNames)

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
