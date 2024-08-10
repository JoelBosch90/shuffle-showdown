package parsers

import (
	"api/lib/helpers"
	"api/lib/wikipedia/models"
	"errors"
	"regexp"
	"strconv"
)

func GetReleaseYear(category models.Category) (uint, error) {
	regex := regexp.MustCompile(`(?i)Catégorie:\s*(Chanson de|Single musical sorti en)\s+(?P<year>\d{1,4})\s*`)
	matches := helpers.GetNamedMatchesForRegex(regex, category.Title)

	if len(matches) == 0 {
		return 0, errors.New("no song category found")
	}

	year, yearError := strconv.ParseUint(matches["year"], 10, 32)
	if yearError != nil {
		return 0, errors.New("error parsing year")
	}

	return uint(year), nil
}
