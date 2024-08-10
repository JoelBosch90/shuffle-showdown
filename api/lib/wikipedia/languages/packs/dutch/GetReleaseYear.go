package dutch

import (
	languageModels "api/lib/wikipedia/languages/models"
	parsers "api/lib/wikipedia/languages/packs/dutch/parsers"
	"api/lib/wikipedia/languages/packs/generic/helpers"
	genericParsers "api/lib/wikipedia/languages/packs/generic/parsers"
	wikipediaModels "api/lib/wikipedia/models"
	"errors"
	"regexp"
	"strings"
)

var releaseYearParsers = []languageModels.ReleaseYearParser{parsers.GetReleaseYearFromSinglesCategory}
var confirmArtistParsers = []languageModels.ConfirmArtistParser{parsers.CategoryConfirmsArtist}
var recognizeRedirectParsers = []languageModels.RecognizeRedirectParser{parsers.IsRedirectCategory}

func isRedirectPage(category wikipediaModels.Category) bool {
	for _, IsRedirect := range recognizeRedirectParsers {
		if IsRedirect(category) {
			return true
		}
	}

	return false
}

func getReleaseYearFromCategory(category wikipediaModels.Category) (uint, error) {
	for _, releaseYearParser := range releaseYearParsers {
		parsedReleaseYear, parseError := releaseYearParser(category)
		if parseError == nil && parsedReleaseYear != 0 {
			return parsedReleaseYear, nil
		}
	}

	return 0, errors.New("no release year found")
}

func categoryConfirmsArtist(category wikipediaModels.Category, artistName string) bool {
	for _, confirmArtistParser := range confirmArtistParsers {
		if confirmArtistParser(category, artistName) {
			return true
		}
	}

	return false
}

func isArtistCategory(category wikipediaModels.Category) bool {
	regex := regexp.MustCompile(`(?i)Categorie:\\s*Nummer van.*`)
	match := regex.FindStringSubmatch(category.Title)

	return len(match) > 0
}

func isDifferentArtistCategory(category wikipediaModels.Category, artistNames []string) bool {
	if !isArtistCategory(category) {
		return false
	}

	for _, artistName := range artistNames {
		regex := regexp.MustCompile("(?i).*" + artistName + ".*")
		match := regex.FindStringSubmatch(category.Title)

		if len(match) > 0 {
			return false
		}
	}

	return true
}

func findOldestReleaseYear(releaseYears []uint) uint {
	oldestReleaseYear := releaseYears[0]

	for _, releaseYear := range releaseYears {
		if releaseYear < oldestReleaseYear {
			oldestReleaseYear = releaseYear
		}
	}

	return oldestReleaseYear
}

/**
 *  1. Loops through the categories of the response.
 *		1. Checks if the category is a redirection page for an early return.
 *		2. Tries to get release year from the category.
 *  	3. Checks if the category is an artist attribution.
 * 	2. Returns an error if no song category is found, or if the artists were not confirmed.
 *  3. Returns the release year otherwise.
 */
func GetReleaseYear(response wikipediaModels.Response, artistNames []string) (uint, error) {
	var releaseYears []uint = []uint{}
	var confirmedArtists = []string{}
	otherArtistsFound := false
	firstPage := response.Query.Pages[0]
	cleanedArtists := helpers.CleanArtistNames(artistNames)

	for _, category := range firstPage.Categories {
		if isRedirectPage(category) {
			return 0, errors.New("redirect page")
		}

		foundReleaseYear, releaseYearError := getReleaseYearFromCategory(category)
		if releaseYearError == nil && foundReleaseYear != 0 {
			releaseYears = append(releaseYears, uint(foundReleaseYear))
		}

		for _, artistName := range cleanedArtists {
			if categoryConfirmsArtist(category, artistName) {
				confirmedArtists = append(confirmedArtists, artistName)
			}
		}

		if isDifferentArtistCategory(category, cleanedArtists) {
			otherArtistsFound = true
		}
	}

	if len(confirmedArtists) == 0 && !genericParsers.TitleConfirmsArtist(firstPage.Title, strings.Join(cleanedArtists, " en ")) {
		return 0, errors.New("artists not confirmed")
	}

	if len(releaseYears) == 0 {
		return 0, errors.New("no release year found")
	}

	if otherArtistsFound && len(releaseYears) > 1 {
		return 0, errors.New("could not distinguish artists/years")
	}

	return findOldestReleaseYear(releaseYears), nil
}
