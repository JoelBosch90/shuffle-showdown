package spotify

import (
	"strconv"
	"strings"
)

func ConvertReleaseDateToIntegers(releaseDate string) (uint, uint, uint) {
	var releaseYear uint = 0
	var releaseMonth uint = 0
	var releaseDay uint = 0

	releaseParts := strings.Split(releaseDate, "-")
	if len(releaseParts) > 0 && len(releaseParts[0]) == 4 {
		convertedYear, conversionError := strconv.ParseUint(releaseParts[0], 10, 32)
		if conversionError == nil {
			releaseYear = uint(convertedYear)
		}
	}
	if len(releaseParts) > 1 && len(releaseParts[1]) == 2 {
		convertedMonth, conversionError := strconv.ParseUint(releaseParts[1], 10, 32)
		if conversionError == nil {
			releaseMonth = uint(convertedMonth)
		}
	}
	if len(releaseParts) > 2 && len(releaseParts[2]) == 2 {
		convertedDay, conversionError := strconv.ParseUint(releaseParts[2], 10, 32)
		if conversionError == nil {
			releaseDay = uint(convertedDay)
		}
	}

	return releaseYear, releaseMonth, releaseDay
}
