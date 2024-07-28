package wikipedia

import (
	wikipediaModels "api/lib/wikipedia/models"
	"errors"
	"strings"
)

type Response struct {
	BatchComplete bool                  `json:"batchcomplete"`
	Query         wikipediaModels.Query `json:"query"`
}

func createTrackTitleSuggestion(trackTitle string, artistNames []string) string {
	if len(artistNames) == 0 {
		return trackTitle + " (song)"
	}

	return trackTitle + " (" + strings.Join(artistNames, " and ") + " song)"
}

func findSuggestionWithoutSuffix(response Response, trackTitle string) string {
	var suggestedTitle string = ""

	if len(response.Query.Pages) == 0 {
		return suggestedTitle
	}

	for _, link := range response.Query.Pages[0].Links {
		suggestedTitle = ParseLinkForCapitalization(link, trackTitle)
		if suggestedTitle != "" && suggestedTitle != trackTitle {
			return suggestedTitle
		}
	}

	return suggestedTitle
}

func findSuggestionWithSongSuffix(response Response) string {
	var suggestedTitle string = ""

	if len(response.Query.Pages) == 0 {
		return suggestedTitle
	}

	for _, link := range response.Query.Pages[0].Links {
		suggestedTitle = ParseLinkForSongSuffix(link)
		if suggestedTitle != "" {
			return suggestedTitle
		}
	}

	return suggestedTitle
}

func findSuggestionWithArtistSuffix(response Response, artistNames []string) string {
	var suggestedTitle string = ""

	if len(response.Query.Pages) == 0 {
		return suggestedTitle
	}

	for _, link := range response.Query.Pages[0].Links {
		suggestedTitle = ParseLinkForArtistsSuffix(link, artistNames)
		if suggestedTitle != "" {
			return suggestedTitle
		}
		suggestedTitle = ParseLinkForArtistsSuffixWithCleaning(link, artistNames)
		if suggestedTitle != "" {
			return suggestedTitle
		}
	}

	return suggestedTitle
}

func trySuggestion(trackTitle string, artistNames []string) (Response, uint) {
	if trackTitle == "" {
		return Response{}, 0
	}

	response, requestError := RequestTrackInfo(trackTitle)
	if requestError != nil {
		return response, 0
	}

	releaseYear, mainCategory, artistsConfirmed := ParseResponseCategories(response, artistNames)
	if mainCategory == "song" && artistsConfirmed {
		return response, releaseYear
	}

	return response, 0
}

func GetTrackReleaseYear(trackTitle string, artistNames []string) (uint, error) {
	trackTitle = CleanTrackTitle(trackTitle)
	originalResponse, releaseYear := trySuggestion(trackTitle, artistNames)
	if releaseYear != 0 {
		return releaseYear, nil
	}

	var suggestion string
	suggestion = findSuggestionWithoutSuffix(originalResponse, trackTitle)
	if suggestion != "" {
		_, releaseYear = trySuggestion(suggestion, artistNames)
	}
	if releaseYear != 0 {
		return releaseYear, nil
	}

	suggestion = findSuggestionWithSongSuffix(originalResponse)
	if suggestion != "" {
		_, releaseYear = trySuggestion(suggestion, artistNames)
	}
	if releaseYear != 0 {
		return releaseYear, nil
	}

	suggestion = findSuggestionWithArtistSuffix(originalResponse, artistNames)
	if suggestion != "" {
		_, releaseYear = trySuggestion(suggestion, artistNames)
	}
	if releaseYear != 0 {
		return releaseYear, nil
	}

	suggestion = createTrackTitleSuggestion(trackTitle, []string{})
	if suggestion != "" {
		_, releaseYear = trySuggestion(suggestion, artistNames)
	}
	if releaseYear != 0 {
		return releaseYear, nil
	}

	suggestion = createTrackTitleSuggestion(trackTitle, artistNames)
	if suggestion != "" {
		_, releaseYear = trySuggestion(suggestion, artistNames)
	}
	if releaseYear != 0 {
		return releaseYear, nil
	}

	return 0, errors.New("not a song page")
}
