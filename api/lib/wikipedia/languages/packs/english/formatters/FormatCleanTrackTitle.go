package formatters

import (
	"api/lib/helpers"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func trimDashedSuffix(trackTitle string) string {
	/**
	 *	^\W*(?P<song>.*)	Captures the song name
	 *	\s+-\s+						Matches the separator between the song name and the suffix
	 *	.*								Matches any suffix
	 */
	regex := regexp.MustCompile(`(?i)^\W*(?P<song>.*)\s+-\s+.*`)
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
	regex := regexp.MustCompile(`(?i)^\W*(?P<song>.*)\s+\(.*\)\s*`)
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

func uncapitalizeArticles(trackTitle string) string {
	/**
	 *	\b(?i)(a|an|the)\b	Matches the articles
	 */
	regex := regexp.MustCompile(`\b(?i)(a|an|the)\b`)
	return regex.ReplaceAllStringFunc(trackTitle, func(match string) string {
		return strings.ToLower(match)
	})
}

func FormatCleanTrackTitle(trackTitle string, _ []string) string {
	withoutDashedSuffix := trimDashedSuffix(trackTitle)
	withoutBracketedSuffix := trimBracketedSuffix(withoutDashedSuffix)
	withoutCommaSuffix := trimCommaSuffix(withoutBracketedSuffix)
	withTrimmedSpaces := strings.TrimSpace(withoutCommaSuffix)
	withTitleCapitalization := cases.Title(language.English).String(withTrimmedSpaces)
	withUncapitalizedArticles := uncapitalizeArticles(withTitleCapitalization)

	return withUncapitalizedArticles
}
