package languages

import (
	"api/lib/wikipedia/languages/models"
	"api/lib/wikipedia/languages/packs/french"
)

var Map = map[models.Language]models.Pack{
	// models.English: english.Pack,
	// models.Dutch:   dutch.Pack,
	// models.German: german.Pack,
	models.French: french.Pack,
}
