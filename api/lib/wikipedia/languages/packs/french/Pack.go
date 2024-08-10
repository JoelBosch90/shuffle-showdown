package french

import (
	"api/lib/wikipedia/languages/models"
)

var Pack = models.Pack{
	Language:            models.French,
	ApiUrl:              "https://fr.wikipedia.org/w/api.php",
	GetLinkSuggestions:  GetLinkSuggestions,
	GetTitleSuggestions: GetTitleSuggestions,
	GetReleaseYear:      GetReleaseYear,
}
