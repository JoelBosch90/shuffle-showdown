package formatters

import (
	"strings"
)

func FormatTrackTitleWithArtistNamesSongSuffix(trackTitle string, artistNames []string) string {
	cleanTitle := FormatCleanTrackTitle(trackTitle, artistNames)

	if len(artistNames) == 0 {
		return cleanTitle + " (song)"
	}

	return cleanTitle + " (" + strings.Join(artistNames, " and ") + " song)"
}
