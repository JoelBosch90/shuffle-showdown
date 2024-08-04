package english

import (
	"api/lib/wikipedia/languages/models"
	"api/lib/wikipedia/languages/packs/english/formatters"
	"api/lib/wikipedia/languages/packs/english/parsers"
)

var Pack = models.Pack{
	Language: models.English,
	ApiUrl:   "https://en.wikipedia.org/w/api.php",
	CategoryParsers: models.CategoryParsers{
		ReleaseYearParsers:       []models.ReleaseYearParser{parsers.GetReleaseYearFromCategory, parsers.GetReleaseYearFromSinglesCategory},
		ConfirmArtistParsers:     []models.ConfirmArtistParser{parsers.CategoryConfirmsArtist, parsers.CategoryConfirmsCleanedArtist},
		RecognizeRedirectParsers: []models.RecognizeRedirectParser{parsers.IsDisambiguationCategory, parsers.IsUnprintworthyRedirectCategory},
	},
	RedirectFormatters: models.RedirectFormatters{
		LinkBasedFormatters:  []models.LinkBasedFormatter{formatters.FindLinkWithOtherCapitalization, formatters.FindLinkWithSongSuffix, formatters.FindLinkWithArtistsSuffix, formatters.FindLinkWithFirstArtistSuffix},
		TitleBasedFormatters: []models.TitleBasedFormatter{formatters.FormatCleanTrackTitle, formatters.FormatTrackTitleWithArtistNamesSongSuffix, formatters.FormatTrackTitleWithOnlySongSuffix, formatters.FormatTrackTitleWithFirstrArtistNameSongSuffix},
	},
}
