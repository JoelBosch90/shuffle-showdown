package wikipedia

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

type Response struct {
	BatchComplete bool  `json:"batchcomplete"`
	Query         Query `json:"query"`
}

func convertMonthStringToInt(monthString string) (uint, error) {
	switch monthString {
	case "January":
		return 1, nil
	case "February":
		return 2, nil
	case "March":
		return 3, nil
	case "April":
		return 4, nil
	case "May":
		return 5, nil
	case "June":
		return 6, nil
	case "July":
		return 7, nil
	case "August":
		return 8, nil
	case "October":
		return 10, nil
	case "November":
		return 11, nil
	case "December":
		return 12, nil
	default:
		return 0, errors.New("error converting month string to int: " + monthString)
	}
}

func getReleaseDate(content string) (uint, uint, uint, error) {
	regex := regexp.MustCompile(`released\s*=\s*(?P<month>\w{3,9})\s*(?P<day>\d{1,2}),\s*(?P<year>\d{1,4})`)
	match := regex.FindStringSubmatch(content)
	if match == nil {
		return 0, 0, 0, errors.New("error parsing release date: " + content)
	}

	month, monthError := convertMonthStringToInt(match[1])
	if monthError != nil {
		return 0, 0, 0, monthError
	}

	day, dayError := strconv.ParseUint(match[2], 10, 32)
	if dayError != nil {
		return 0, 0, 0, dayError
	}

	year, yearError := strconv.ParseUint(match[3], 10, 32)
	if yearError != nil {
		return 0, 0, 0, yearError
	}

	return uint(year), month, uint(day), nil
}

func GetTrackReleaseDate(trackName string, artistNames []string) (uint, uint, uint, error) {
	headers := []Header{}
	params := []Param{
		{Name: "action", Value: "query"},
		{Name: "prop", Value: "revisions"},
		{Name: "rvprop", Value: "content"},
		{Name: "titles", Value: trackName},
		{Name: "utf8", Value: "1"},
		{Name: "rvsection", Value: "0"},
		{Name: "rvslots", Value: "*"},
		{Name: "format", Value: "json"},
		{Name: "formatversion", Value: "2"},
	}

	response, requestError := ApiRequest(http.MethodGet, headers, params)
	if requestError != nil {
		return 0, 0, 0, requestError
	}

	var formattedResponse Response
	decoder := json.NewDecoder(response.Body)
	decodeError := decoder.Decode(&formattedResponse)
	if decodeError != nil {
		return 0, 0, 0, decodeError
	}

	revisions := formattedResponse.Query.Pages[0].Revisions
	if len(revisions) == 0 {
		return 0, 0, 0, errors.New("no revisions found")
	}

	lastRevision := revisions[len(revisions)-1]
	content := lastRevision.Slots.Main.Content

	releaseYear, releaseMonth, releaseDay, parseError := getReleaseDate(content)
	if parseError != nil {
		return 0, 0, 0, parseError
	}

	return releaseYear, releaseMonth, releaseDay, nil
}
