package wikipedia_languages

import (
	english "api/lib/wikipedia/languages/english"
	wikipediaModels "api/lib/wikipedia/models"
)

type Language string

const (
	English Language = "en"
	Dutch   Language = "nl"
	German  Language = "de"
	French  Language = "fr"
	Spanish Language = "es"
)

type ReleaseYearParser func(category wikipediaModels.Category) (uint, error)
type ConfirmArtistParser func(category wikipediaModels.Category, artistName string) bool
type RedirectionParser func(category wikipediaModels.Category) bool

type CategoryParsers struct {
	ReleaseYearParsers   []ReleaseYearParser
	ConfirmArtistParsers []ConfirmArtistParser
	RedirectParsers      []RedirectionParser
}

type LanguagePack struct {
	Language        Language        `json:"language"`
	ApiUrl          string          `json:"apiUrl"`
	CategoryParsers CategoryParsers `json:"categoryParsers"`
}

var LanguageMap = map[Language]LanguagePack{
	English: {
		Language: English,
		ApiUrl:   "https://en.wikipedia.org/w/api.php",
		CategoryParsers: CategoryParsers{
			ReleaseYearParsers:   []ReleaseYearParser{english.GetReleaseYearFromCategory},
			ConfirmArtistParsers: []ConfirmArtistParser{english.CategoryConfirmsArtist, english.CategoryConfirmsCleanedArtist},
			RedirectParsers:      []RedirectionParser{english.IsDisambiguationCategory, english.IsUnprintworthyRedirectCategory},
		},
	},
	// Dutch: {
	// 	Language: Dutch,
	// 	ApiUrl:   "https://nl.wikipedia.org/w/api.php",
	// },
	// German: {
	// 	Language: German,
	// 	ApiUrl:   "https://de.wikipedia.org/w/api.php",
	// },
	// French: {
	// 	Language: French,
	// 	ApiUrl:   "https://fr.wikipedia.org/w/api.php",
	// },
	// Spanish: {
	// 	Language: Spanish,
	// 	ApiUrl:   "https://es.wikipedia.org/w/api.php",
	// },
}
