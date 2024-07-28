package wikipedia

import (
	helpers "api/lib/helpers"
	"regexp"
	"strings"
)

func CleanTrackTitle(trackTitle string) string {
	/**
	 *	^\W*(?P<song>.*)				Captures the song name
	 *	\s+-\s*									Matches the separator between the song name and the remaster/mix
	 *	\d{1,4}\s*remaster			Matches a remaster suffix
	 *	|.*\WMix								Matches a mix suffix
	 */
	regex := regexp.MustCompile(`(?i)^\W*(?P<song>.*)\s+-\s*(?:\d{1,4}\s*remaster|.*\WMix).*`)
	matches := helpers.GetNamedMatchesForRegex(regex, trackTitle)

	if len(matches) != 0 {
		return strings.TrimSpace(matches["song"])
	}

	return strings.TrimSpace(trackTitle)
}
