package english

import (
	"api/lib/wikipedia/languages/models"
)

var Pack = models.Pack{
	Language:            models.English,
	ApiUrl:              "https://en.wikipedia.org/w/api.php",
	GetLinkSuggestions:  GetLinkSuggestions,
	GetTitleSuggestions: GetTitleSuggestions,
	GetReleaseYear:      GetReleaseYear,
}
