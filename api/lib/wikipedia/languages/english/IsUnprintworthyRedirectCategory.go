package wikipedia_languages_english

import (
	wikipediaModels "api/lib/wikipedia/models"
	"regexp"
)

func IsUnprintworthyRedirectCategory(category wikipediaModels.Category) bool {
	regex := regexp.MustCompile(`(?i)Category:\s*Unprintworthy\s+redirects`)
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
