package wikipedia_languages_english

import (
	wikipediaModels "api/lib/wikipedia/models"
	"regexp"
)

func IsDisambiguationCategory(category wikipediaModels.Category) bool {
	regex := regexp.MustCompile(`(?i)Category:\s*Disambiguation\s+pages`)
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
