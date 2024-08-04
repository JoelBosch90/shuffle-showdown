package parsers

import (
	"api/lib/wikipedia/models"
	"regexp"
)

func IsDisambiguationCategory(category models.Category) bool {
	regex := regexp.MustCompile(`(?i)Category:\s*Disambiguation\s+pages`)
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
