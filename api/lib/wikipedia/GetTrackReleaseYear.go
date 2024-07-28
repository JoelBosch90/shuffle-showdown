package wikipedia

import (
	wikipediaModels "api/lib/wikipedia/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Response struct {
	BatchComplete bool                  `json:"batchcomplete"`
	Query         wikipediaModels.Query `json:"query"`
}

func parseSongCategory(category wikipediaModels.Category) (bool, uint) {
	regex := regexp.MustCompile(`(?i)Category:\s*(?P<releaseYear>\d{1,4})\s*songs`)
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

func parseArtistCategory(category wikipediaModels.Category, artistName string) bool {
	regex := regexp.MustCompile("(?i)Category:\\s*" + artistName + ".*songs")
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

func parsePage(page wikipediaModels.Page, artistNames []string) (bool, uint) {
	var releaseYear uint
	var hasSongCategory bool
	var artistsInCategories []string

	for _, category := range page.Categories {
		isSongCategory, songReleaseYear := parseSongCategory(category)

		if isSongCategory {
			hasSongCategory = true
			releaseYear = songReleaseYear
		}

		for _, artistName := range artistNames {
			if parseArtistCategory(category, artistName) {
				artistsInCategories = append(artistsInCategories, artistName)
			}
		}
	}

	log.Println("ARTISTS IN CATEGORIES", artistsInCategories)
	log.Println("ARTIST NAMES", artistNames)
	log.Println("HAS SONG CATEGORY", hasSongCategory)
	log.Println("INCLUDES STRINGS", includesStrings(artistsInCategories, artistNames))
	log.Println("RELEASE YEAR", releaseYear)

	return hasSongCategory && includesStrings(artistsInCategories, artistNames), releaseYear
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

func GetTrackReleaseYear(trackTitle string, artistNames []string) (uint, error) {
	headers := []Header{}
	params := []Param{
		{Name: "action", Value: "query"},
		{Name: "prop", Value: "categories"},
		{Name: "titles", Value: cleanUpTrackTitle(trackTitle)},
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

	page := formattedResponse.Query.Pages[0]
	log.Println("ARTIST NAMES", artistNames)
	isSongPage, releaseYear := parsePage(page, artistNames)
	log.Println("TRACK NAME", trackTitle, cleanUpTrackTitle(trackTitle))
	if !isSongPage {
		log.Println("NOT A SONG PAGE?", page)
		return 0, errors.New("not a song page")
	}

	return releaseYear, nil
}
