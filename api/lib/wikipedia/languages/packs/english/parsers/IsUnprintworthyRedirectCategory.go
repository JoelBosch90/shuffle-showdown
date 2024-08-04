package parsers

import (
	"api/lib/wikipedia/models"
	"regexp"
)

func IsUnprintworthyRedirectCategory(category models.Category) bool {
	regex := regexp.MustCompile(`(?i)Category:\s*Unprintworthy\s+redirects`)
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
