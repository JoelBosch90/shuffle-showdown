package spotify

import (
	"strconv"
	"strings"
)

func ConvertReleaseDateToReleaseYear(releaseDate string) uint {
	var releaseYear uint = 0

	releaseParts := strings.Split(releaseDate, "-")
	if len(releaseParts) > 0 && len(releaseParts[0]) == 4 {
		convertedYear, conversionError := strconv.ParseUint(releaseParts[0], 10, 32)
		if conversionError == nil {
			releaseYear = uint(convertedYear)
		}
	}

	return releaseYear
}
