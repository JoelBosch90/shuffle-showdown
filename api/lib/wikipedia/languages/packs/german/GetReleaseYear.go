package german

import (
	languageModels "api/lib/wikipedia/languages/models"
	parsers "api/lib/wikipedia/languages/packs/german/parsers"
	wikipediaModels "api/lib/wikipedia/models"
	"errors"
	"log"
)

var releaseYearParsers = []languageModels.ReleaseYearParser{parsers.GetReleaseYearFromSinglesCategory}
var confirmArtistParsers = []languageModels.ConfirmArtistParser{parsers.CategoryConfirmsArtist}
var recognizeRedirectParsers = []languageModels.RecognizeRedirectParser{}

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

/**
 *  1. Loops through the categories of the response.
 *		1. Checks if the category is a redirection page for an early return.
 *		2. Tries to get release year from the category.
 *  	3. Checks if the category is an artist attribution.
 * 	2. Returns an error if no song category is found, or if the artists were not confirmed.
 *  3. Returns the release year otherwise.
 */
func GetReleaseYear(response wikipediaModels.Response, artistNames []string) (uint, error) {
	var releaseYear uint = uint(0)
	var confirmedArtists = []string{}

	for _, category := range response.Query.Pages[0].Categories {
		if isRedirectPage(category) {
			log.Println("IS REDIRECT PAGE")
			return 0, errors.New("redirect page")
		}

		foundReleaseYear, releaseYearError := getReleaseYearFromCategory(category)
		log.Println("FOUND RELEASE YEAR", foundReleaseYear, releaseYearError)
		if releaseYearError == nil && foundReleaseYear != 0 {
			if releaseYear == 0 {
				log.Println("RELEASE YEAR FOUND", foundReleaseYear)
				releaseYear = uint(foundReleaseYear)
			} else {
				log.Println("MULTIPLE RELEASE YEARS FOUND")
				return 0, errors.New("multiple release years found")
			}
		}

		for _, artistName := range artistNames {
			if categoryConfirmsArtist(category, artistName) {
				log.Println("CONFIRMED ARTIST", artistName)
				confirmedArtists = append(confirmedArtists, artistName)
			}
		}
	}

	if len(confirmedArtists) == 0 {
		log.Println("ARTISTS NOT CONFIRMED")
		return 0, errors.New("artists not confirmed")
	}

	if releaseYear == 0 {
		log.Println("NO RELEASE YEAR FOUND")
		return 0, errors.New("no release year found")
	}

	return releaseYear, nil
}
