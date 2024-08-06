package dutch

import (
	"api/lib/wikipedia/languages/models"
)

var Pack = models.Pack{
	Language:           models.Dutch,
	ApiUrl:             "https://nl.wikipedia.org/w/api.php",
	CategoryParsers:    models.CategoryParsers{},
	RedirectFormatters: models.RedirectFormatters{},
}
