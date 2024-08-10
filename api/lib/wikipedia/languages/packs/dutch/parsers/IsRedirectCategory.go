package parsers

import (
	"api/lib/wikipedia/models"
	"regexp"
)

func IsRedirectCategory(category models.Category) bool {
	regex := regexp.MustCompile(`(?i)Categorie:\s*Wikipedia:\s*Doorverwijspagina`)
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
