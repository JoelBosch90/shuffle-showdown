package dutch

import (
	"api/lib/wikipedia/languages/models"
)

var Pack = models.Pack{
	Language:            models.Dutch,
	ApiUrl:              "https://nl.wikipedia.org/w/api.php",
	GetLinkSuggestions:  GetLinkSuggestions,
	GetTitleSuggestions: GetTitleSuggestions,
	GetReleaseYear:      GetReleaseYear,
}
