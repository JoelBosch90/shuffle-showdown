package formatters

import "api/lib/helpers"

func FormatCleanTrackTitleWithoutCapitalization(trackTitle string, artistNames []string) string {
	withCapitalization := FormatCleanTrackTitle(trackTitle, artistNames)

	return helpers.ConvertToLowerCaseExceptFirstCharacter(withCapitalization)
}
