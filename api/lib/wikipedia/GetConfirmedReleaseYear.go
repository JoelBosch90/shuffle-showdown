package wikipedia

import (
	languageModels "api/lib/wikipedia/languages/models"
	wikipediaModels "api/lib/wikipedia/models"
	"errors"
)

func isRedirectPage(categoryParsers languageModels.CategoryParsers, category wikipediaModels.Category) bool {
	for _, IsRedirect := range categoryParsers.RecognizeRedirectParsers {
		if IsRedirect(category) {
			return true
		}
	}

	return false
}

func getReleaseYearFromCategory(categoryParsers languageModels.CategoryParsers, category wikipediaModels.Category) (uint, error) {
	for _, releaseYearParser := range categoryParsers.ReleaseYearParsers {
		parsedReleaseYear, parseError := releaseYearParser(category)
		if parseError == nil && parsedReleaseYear != 0 {
			return parsedReleaseYear, nil
		}
	}

	return 0, errors.New("no release year found")
}

func categoryConfirmsArtist(categoryParsers languageModels.CategoryParsers, category wikipediaModels.Category, artistName string) bool {
	for _, confirmArtistParser := range categoryParsers.ConfirmArtistParsers {
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
func GetConfirmedReleaseYear(languagePack languageModels.Pack, response wikipediaModels.Response, artistNames []string) (uint, error) {
	var releaseYear uint = uint(0)
	var confirmedArtists = []string{}
	categoryParsers := languagePack.CategoryParsers

	for _, category := range response.Query.Pages[0].Categories {
		if isRedirectPage(categoryParsers, category) {
			return 0, errors.New("redirect page")
		}

		foundReleaseYear, releaseYearError := getReleaseYearFromCategory(categoryParsers, category)
		if releaseYearError == nil {
			releaseYear = uint(foundReleaseYear)
		}

		for _, artistName := range artistNames {
			if categoryConfirmsArtist(categoryParsers, category, artistName) {
				confirmedArtists = append(confirmedArtists, artistName)
			}
		}
	}

	if len(confirmedArtists) == 0 {
		return 0, errors.New("artists not confirmed")
	}

	if releaseYear == 0 {
		return 0, errors.New("no release year found")
	}

	return releaseYear, nil
}
