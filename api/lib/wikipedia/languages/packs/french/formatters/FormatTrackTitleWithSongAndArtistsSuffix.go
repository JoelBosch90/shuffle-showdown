package formatters

import (
	"regexp"
	"strings"
)

func FormatTrackTitleWithSongAndArtistsSuffix(trackTitle string, artistNames []string) string {
	return FormatCleanTrackTitle(trackTitle, []string{}) + " (chanson de " + regexp.QuoteMeta(strings.Join(artistNames, " et ")) + ")"
}
