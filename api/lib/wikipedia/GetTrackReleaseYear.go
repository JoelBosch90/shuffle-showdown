package wikipedia

import (
	wikipediaModels "api/lib/wikipedia/models"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Response struct {
	BatchComplete bool                  `json:"batchcomplete"`
	Query         wikipediaModels.Query `json:"query"`
}

func parseSongCategory(category wikipediaModels.Category) (uint, bool) {
	regex := regexp.MustCompile(`(?i)Category:\s*(?P<releaseYear>\d{1,4})\s+songs`)
	match := regex.FindStringSubmatch(category.Title)

	if match == nil {
		return 0, false
	}

	year, yearError := strconv.ParseUint(match[1], 10, 32)
	if yearError != nil {
		return 0, false
	}

	return uint(year), true
}

func parseDisambiguationCategory(category wikipediaModels.Category) bool {
	regex := regexp.MustCompile(`(?i)Category:\s*Disambiguation\s+pages`)
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}

func parseUnprintworthyRedirectCategory(category wikipediaModels.Category) bool {
	regex := regexp.MustCompile(`(?i)Category:\s*Unprintworthy\s+redirects`)
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}

func parseArtistCategory(category wikipediaModels.Category, artistName string) bool {
	regex := regexp.MustCompile("(?i)Category:\\s*" + regexp.QuoteMeta(cleanUpArtistName(artistName)) + "(\\s+\\((band|group)\\))?\\s+songs")
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}

func includesString(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}

	return false
}

func includesStrings(haystack []string, needles []string) bool {
	for _, needle := range needles {
		if !includesString(haystack, needle) {
			return false
		}
	}

	return true
}

func parseResponseCategories(response Response, artistNames []string) (uint, string, bool) {
	var releaseYear uint = uint(0)
	var mainCategory string = ""
	var artistsInCategories []string = []string{}

	if len(response.Query.Pages) == 0 {
		return releaseYear, mainCategory, false
	}

	page := response.Query.Pages[0]
	for _, category := range page.Categories {
		songReleaseYear, isSongCategory := parseSongCategory(category)

		if parseDisambiguationCategory(category) {
			mainCategory = "disambiguation"
			break
		}

		if parseUnprintworthyRedirectCategory(category) {
			mainCategory = "unprintworthy"
			break
		}

		if isSongCategory {
			mainCategory = "song"
			releaseYear = songReleaseYear
		}

		for _, artistName := range artistNames {
			if parseArtistCategory(category, artistName) {
				artistsInCategories = append(artistsInCategories, artistName)
			}
		}
	}

	return releaseYear, mainCategory, includesStrings(artistsInCategories, artistNames)
}

func parseLink(link wikipediaModels.Link, artistNames []string) string {
	var regex *regexp.Regexp
	if len(artistNames) == 0 {
		regex = regexp.MustCompile(`(?i).*\(song\)`)
	} else {
		regex = regexp.MustCompile("(?i).*\\(" + regexp.QuoteMeta(strings.Join(cleanUpArtistNames(artistNames), " and ")) + " song\\)")
	}
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}

func parseDifferentCapitalization(link wikipediaModels.Link, trackTitle string) string {
	regex := regexp.MustCompile("(?i)^\\s*" + regexp.QuoteMeta(cleanUpTrackTitle(trackTitle)) + "\\s*$")
	match := regex.FindStringSubmatch(link.Title)

	if match == nil {
		return ""
	}

	return match[0]
}

func findTrackTitleSuggestion(response Response, trackTitle string, artistNames []string) string {
	var suggestedTitle string = ""

	if len(response.Query.Pages) == 0 {
		return suggestedTitle
	}

	for _, link := range response.Query.Pages[0].Links {
		suggestedTitle = parseDifferentCapitalization(link, trackTitle)
		if suggestedTitle != "" && suggestedTitle != trackTitle {
			return suggestedTitle
		}
	}

	for _, link := range response.Query.Pages[0].Links {
		suggestedTitle = parseLink(link, []string{})
		if suggestedTitle != "" {
			return suggestedTitle
		}
	}

	for _, link := range response.Query.Pages[0].Links {
		suggestedTitle = parseLink(link, artistNames)
		if suggestedTitle != "" {
			return suggestedTitle
		}
	}

	return suggestedTitle
}

func createTrackTitleSuggestion(trackTitle string, artistNames []string) string {
	if len(artistNames) == 0 {
		return trackTitle + " (song)"
	}

	return trackTitle + " (" + strings.Join(artistNames, " and ") + " song)"
}

func cleanUpArtistName(artistName string) string {
	regex := regexp.MustCompile(`(?i)The\s+(?P<mainName>.*)`)
	match := regex.FindStringSubmatch(artistName)

	if match != nil {
		return match[1]
	}

	return artistName
}

func cleanUpArtistNames(artistNames []string) []string {
	var cleanedArtistNames []string

	for _, artistName := range artistNames {
		cleanedArtistNames = append(cleanedArtistNames, cleanUpArtistName(artistName))
	}

	return cleanedArtistNames
}

func cleanUpTrackTitle(trackTitle string) string {
	/**
	 *	^\W*(?P<songName>.*)		Captures the song name
	 *	\s+-\s*									Matches the separator between the song name and the remaster/mix
	 *	\d{1,4}\s*remaster			Matches a remaster suffix
	 *	|.*\WMix								Matches a mix suffix
	 */
	regex := regexp.MustCompile(`(?i)^\W*(?P<songName>.*)\s+-\s*(?:\d{1,4}\s*remaster|.*\WMix).*`)
	match := regex.FindStringSubmatch(trackTitle)

	if match != nil {
		return match[1]
	}

	return trackTitle
}

func makeRequest(trackTitle string) (Response, error) {
	var formattedResponse Response

	headers := []Header{}
	params := []Param{
		{Name: "action", Value: "query"},
		{Name: "prop", Value: "categories|links"},
		{Name: "titles", Value: cleanUpTrackTitle(trackTitle)},
		{Name: "utf8", Value: "1"},
		{Name: "format", Value: "json"},
		{Name: "formatversion", Value: "2"},
		{Name: "cllimit", Value: "max"},
		{Name: "pllimit", Value: "max"},
	}

	response, requestError := ApiRequest(http.MethodGet, headers, params)
	if requestError != nil {
		return formattedResponse, requestError
	}

	decoder := json.NewDecoder(response.Body)
	decodeError := decoder.Decode(&formattedResponse)
	if decodeError != nil {
		return formattedResponse, decodeError
	}

	return formattedResponse, nil
}

func GetTrackReleaseYear(trackTitle string, artistNames []string) (uint, error) {
	response, requestError := makeRequest(trackTitle)
	if requestError != nil {
		return 0, requestError
	}

	releaseYear, mainCategory, artistsConfirmed := parseResponseCategories(response, artistNames)

	if mainCategory != "song" || !artistsConfirmed {
		titleSuggestion := findTrackTitleSuggestion(response, trackTitle, artistNames)
		if titleSuggestion == "" {
			titleSuggestion = createTrackTitleSuggestion(trackTitle, []string{})
		}

		response, requestError = makeRequest(titleSuggestion)
		if requestError != nil {
			return 0, requestError
		}
		releaseYear, mainCategory, artistsConfirmed = parseResponseCategories(response, artistNames)
	}

	if mainCategory != "song" || !artistsConfirmed {
		titleSuggestion := createTrackTitleSuggestion(trackTitle, artistNames)

		response, requestError = makeRequest(titleSuggestion)
		if requestError != nil {
			return 0, requestError
		}

		releaseYear, mainCategory, artistsConfirmed = parseResponseCategories(response, artistNames)
	}

	if mainCategory != "song" || !artistsConfirmed {
		return 0, errors.New("not a song page")
	}

	return releaseYear, nil
}
