package languages

import (
	"api/lib/wikipedia/languages/models"
	"api/lib/wikipedia/languages/packs/dutch"
	"api/lib/wikipedia/languages/packs/english"
	"api/lib/wikipedia/languages/packs/french"
	"api/lib/wikipedia/languages/packs/german"
)

var Map = map[models.Language]models.Pack{
	models.English: english.Pack,
	models.Dutch:   dutch.Pack,
	models.German:  german.Pack,
	models.French:  french.Pack,
}
