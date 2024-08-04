package formatters

import (
	"api/lib/helpers"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func trimSuffix(trackTitle string) string {
	/**
	 *	^\W*(?P<song>.*)				Captures the song name
	 *	\s+-\s*									Matches the separator between the song name and the remaster/mix
	 *	\d{1,4}\s*							Matches a year
	 *	(?:re)?(?:master|mix)		Matches a (re)master or (re)mix suffix
	 *	.*											Matches any remaining characters
	 */
	regex := regexp.MustCompile(`(?i)^\W*(?P<song>.*)\s+-\s*\d{0,4}\s*(?:re)?(?:master|mix).*`)
	matches := helpers.GetNamedMatchesForRegex(regex, trackTitle)

	if len(matches) != 0 {
		return strings.TrimSpace(matches["song"])
	}

	return trackTitle
}

func FormatCleanTrackTitle(trackTitle string, _ []string) string {
	withoutSuffix := trimSuffix(trackTitle)
	withTrimmedSpaces := strings.TrimSpace(withoutSuffix)
	withTitleCapitalization := cases.Title(language.English).String(withTrimmedSpaces)

	return withTitleCapitalization
}
