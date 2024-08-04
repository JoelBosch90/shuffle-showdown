package languages

import (
	"api/lib/wikipedia/languages/models"
	"api/lib/wikipedia/languages/packs/english"
)

var Map = map[models.Language]models.Pack{
	models.English: english.Pack,
	// languageModels.Dutch: {
	// 	Language: languageModels.Dutch,
	// 	ApiUrl:   "https://nl.wikipedia.org/w/api.php",
	// },
	// languageModels.German: {
	// 	Language: languageModels.German,
	// 	ApiUrl:   "https://de.wikipedia.org/w/api.php",
	// },
	// languageModels.French: {
	// 	Language: languageModels.French,
	// 	ApiUrl:   "https://fr.wikipedia.org/w/api.php",
	// },
	// languageModels.Spanish: {
	// 	Language: languageModels.Spanish,
	// 	ApiUrl:   "https://es.wikipedia.org/w/api.php",
	// },
}
