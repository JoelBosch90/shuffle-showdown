package languages

import (
	"api/lib/wikipedia/languages/models"
	"api/lib/wikipedia/languages/packs/german"
)

var Map = map[models.Language]models.Pack{
	// models.English: english.Pack,
	// models.Dutch:   dutch.Pack,
	models.German: german.Pack,
	// models.French: {
	// 	Language: models.French,
	// 	ApiUrl:   "https://fr.wikipedia.org/w/api.php",
	// },
	// models.Spanish: {
	// 	Language: models.Spanish,
	// 	ApiUrl:   "https://es.wikipedia.org/w/api.php",
	// },
}
