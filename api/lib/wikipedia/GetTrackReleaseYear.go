package wikipedia

import (
	languages "api/lib/wikipedia/languages"
	languageModels "api/lib/wikipedia/languages/models"
	wikipediaModels "api/lib/wikipedia/models"
	"errors"
	"log"
)

func getTitleSuggestions(languagePack languageModels.Pack, trackTitle string, artistNames []string) []string {
	var titleSuggestions []string = []string{trackTitle}
	for _, formatter := range languagePack.RedirectFormatters.TitleBasedFormatters {
		titleSuggestions = append(titleSuggestions, formatter(trackTitle, artistNames))
	}

	return titleSuggestions
}

func getLinkSuggestions(languagePack languageModels.Pack, page wikipediaModels.Page, trackTitle string, artistNames []string) []string {
	var linkSuggestions []string = []string{trackTitle}
	for _, formatter := range languagePack.RedirectFormatters.LinkBasedFormatters {
		linkSuggestions = append(linkSuggestions, formatter(page, trackTitle, artistNames))
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
	trackTitle = "Wish You were here"
	artistNames = []string{"Incubus"}

	for _, languagePack := range languageMaps {
		titleSuggestions := getTitleSuggestions(languagePack, trackTitle, artistNames)
		log.Println("TITLE SUGGESTIONS", titleSuggestions)

		for _, suggestion := range titleSuggestions {
			log.Println("PROCESSING TITLE SUGGESTION: ", suggestion)
			response, releaseYear := trySuggestion(languagePack, suggestion, artistNames)
			log.Println("PROCESSED TITLE SUGGESTION: ", response, releaseYear)

			if releaseYear != 0 {
				return releaseYear, nil
			}

			if isPageMissing(response) {
				continue
			}

			firstPage := response.Query.Pages[0]
			linkSuggestions := getLinkSuggestions(languagePack, firstPage, suggestion, artistNames)
			log.Println("LINK SUGGESTIONS", linkSuggestions)

			for _, suggestion := range linkSuggestions {
				log.Println("PROCESSING LINK SUGGESTION: ", suggestion)
				response, releaseYear = trySuggestion(languagePack, suggestion, artistNames)
				log.Println("PROCESSED LINK SUGGESTION: ", response, releaseYear)

				if releaseYear != 0 {
					return releaseYear, nil
				}
			}
		}
	}

	return 0, errors.New("not a song page")
}
