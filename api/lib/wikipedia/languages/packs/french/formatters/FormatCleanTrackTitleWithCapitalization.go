package formatters

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func uncapitalizeArticles(trackTitle string) string {
	/**
	 *	\b(?i)(le|la|l'|un|une)\b	Matches the articles
	 */
	regex := regexp.MustCompile(`\b(?i)(le|la|l'|un|une)\b`)
	return regex.ReplaceAllStringFunc(trackTitle, func(match string) string {
		return strings.ToLower(match)
	})
}

func FormatCleanTrackTitleWithCapitalization(trackTitle string, artistNames []string) string {
	withoutTitleCapitalization := FormatCleanTrackTitle(trackTitle, artistNames)
	withTitleCapitalization := cases.Title(language.French).String(withoutTitleCapitalization)
	withUncapitalizedArticles := uncapitalizeArticles(withTitleCapitalization)

	return withUncapitalizedArticles
}
