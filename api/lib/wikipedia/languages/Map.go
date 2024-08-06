package languages

import (
	"api/lib/wikipedia/languages/models"
	"api/lib/wikipedia/languages/packs/dutch"
	"api/lib/wikipedia/languages/packs/english"
)

var Map = map[models.Language]models.Pack{
	models.English: english.Pack,
	models.Dutch:   dutch.Pack,
	// models.German: {
	// 	Language: models.German,
	// 	ApiUrl:   "https://de.wikipedia.org/w/api.php",
	// },
	// models.French: {
	// 	Language: models.French,
	// 	ApiUrl:   "https://fr.wikipedia.org/w/api.php",
	// },
	// models.Spanish: {
	// 	Language: models.Spanish,
	// 	ApiUrl:   "https://es.wikipedia.org/w/api.php",
	// },
}
