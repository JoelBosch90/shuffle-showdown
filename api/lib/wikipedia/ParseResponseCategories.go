package wikipedia

import (
	helpers "api/lib/helpers"
	wikipediaModels "api/lib/wikipedia/models"
	"regexp"
	"strconv"
)

func parseSongCategory(category wikipediaModels.Category) (uint, bool) {
	regex := regexp.MustCompile(`(?i)Category:\s*(?P<year>\d{1,4})\s+(songs|ballads)`)
	matches := helpers.GetNamedMatchesForRegex(regex, category.Title)

	if len(matches) == 0 {
		return 0, false
	}

	year, yearError := strconv.ParseUint(matches["year"], 10, 32)
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
	regex := regexp.MustCompile("(?i)Category:\\s*" + regexp.QuoteMeta(artistName) + "(\\s+\\((band|group)\\))?\\s+songs")
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}

func parseArtistCategoryWithCleaning(category wikipediaModels.Category, artistName string) bool {
	return parseArtistCategory(category, CleanArtistName(artistName))
}

func ParseResponseCategories(response Response, artistNames []string) (uint, string, bool) {
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
			if parseArtistCategoryWithCleaning(category, artistName) {
				artistsInCategories = append(artistsInCategories, artistName)
			}
		}
	}

	return releaseYear, mainCategory, helpers.IncludesStrings(artistsInCategories, artistNames)
}
