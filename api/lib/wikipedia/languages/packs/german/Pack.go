package german

import (
	"api/lib/wikipedia/languages/models"
)

var Pack = models.Pack{
	Language:            models.German,
	ApiUrl:              "https://de.wikipedia.org/w/api.php",
	GetLinkSuggestions:  GetLinkSuggestions,
	GetTitleSuggestions: GetTitleSuggestions,
	GetReleaseYear:      GetReleaseYear,
}
