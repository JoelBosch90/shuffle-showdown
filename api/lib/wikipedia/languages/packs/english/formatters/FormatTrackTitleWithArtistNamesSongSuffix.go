package formatters

import (
	"strings"
)

func FormatTrackTitleWithArtistNamesSongSuffix(trackTitle string, artistNames []string) string {
	if len(artistNames) == 0 {
		return trackTitle + " (song)"
	}

	return trackTitle + " (" + strings.Join(artistNames, " and ") + " song)"
}
