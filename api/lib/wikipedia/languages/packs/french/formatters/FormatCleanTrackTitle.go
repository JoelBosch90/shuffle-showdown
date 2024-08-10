package formatters

import (
	"api/lib/helpers"
	"regexp"
	"strings"
)

func trimDashedSuffix(trackTitle string) string {
	/**
	 *	^\W*(?P<song>.*)	Captures the song name
	 *	\s+-\s+						Matches the separator between the song name and the suffix
	 *	.*								Matches any suffix
	 */
	regex := regexp.MustCompile(`(?i)^(?P<song>.*?)\s+-\s+.*`)
	matches := helpers.GetNamedMatchesForRegex(regex, trackTitle)

	if len(matches) != 0 {
		return strings.TrimSpace(matches["song"])
	}

	return trackTitle
}

func trimBracketedSuffix(trackTitle string) string {
	/**
	 *	^\W*(?P<song>.*)	Captures the song name
	 *	\s+\(.*\)$			  Matches the separator around the bracketed suffix
	 *	.*								Matches any suffix
	 */
	regex := regexp.MustCompile(`(?i)^(?P<song>.*)\s+\(.*\)\s*`)
	matches := helpers.GetNamedMatchesForRegex(regex, trackTitle)

	if len(matches) != 0 {
		return strings.TrimSpace(matches["song"])
	}

	return trackTitle
}

func trimCommaSuffix(trackTitle string) string {
	/**
	 *	^\W*(?P<song>.*)	Captures the song name
	 *	\s+,\s+						Matches the separator between the song name and the suffix
	 *	.*								Matches any suffix
	 */
	regex := regexp.MustCompile(`(?i)^\W*(?P<song>.*)\s*,.*`)
	matches := helpers.GetNamedMatchesForRegex(regex, trackTitle)

	if len(matches) != 0 {
		return strings.TrimSpace(matches["song"])
	}

	return trackTitle
}

func FormatCleanTrackTitle(trackTitle string, _ []string) string {
	withoutDashedSuffix := trimDashedSuffix(trackTitle)
	withoutBracketedSuffix := trimBracketedSuffix(withoutDashedSuffix)
	withoutCommaSuffix := trimCommaSuffix(withoutBracketedSuffix)
	withTrimmedSpaces := strings.TrimSpace(withoutCommaSuffix)

	return withTrimmedSpaces
}
