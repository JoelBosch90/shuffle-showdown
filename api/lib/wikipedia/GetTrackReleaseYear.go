package wikipedia

import (
	wikipediaModels "api/lib/wikipedia/models"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

type Response struct {
	BatchComplete bool                  `json:"batchcomplete"`
	Query         wikipediaModels.Query `json:"query"`
}

func parseSongCategory(category wikipediaModels.Category) (bool, uint) {
	regex := regexp.MustCompile(`Category:\d{1-4} songs`)
	match := regex.FindStringSubmatch(category.Title)

	if match == nil {
		return false, 0
	}

	year, yearError := strconv.ParseUint(match[1], 10, 32)
	if yearError != nil {
		return false, 0
	}

	return true, uint(year)
}

func parsePage(page wikipediaModels.Page) (bool, uint) {
	for _, category := range page.Categories {
		isSongCategory, releaseYear := parseSongCategory(category)

		if isSongCategory {
			return true, releaseYear
		}
	}

	return false, 0
}

func GetTrackReleaseYear(trackName string, artistNames []string) (uint, error) {
	headers := []Header{}
	params := []Param{
		{Name: "action", Value: "query"},
		{Name: "prop", Value: "revisions|categories"},
		{Name: "titles", Value: trackName},
		{Name: "utf8", Value: "1"},
		{Name: "format", Value: "json"},
		{Name: "formatversion", Value: "2"},
		{Name: "cllimit", Value: "max"},
	}

	response, requestError := ApiRequest(http.MethodGet, headers, params)
	if requestError != nil {
		return 0, requestError
	}

	var formattedResponse Response
	decoder := json.NewDecoder(response.Body)
	decodeError := decoder.Decode(&formattedResponse)
	if decodeError != nil {
		return 0, decodeError
	}

	revisions := formattedResponse.Query.Pages[0].Revisions
	if len(revisions) == 0 {
		return 0, errors.New("no revisions found")
	}

	page := formattedResponse.Query.Pages[0]
	isSongPage, releaseYear := parsePage(page)
	if !isSongPage {
		return 0, errors.New("not a song page")
	}

	return releaseYear, nil
}
