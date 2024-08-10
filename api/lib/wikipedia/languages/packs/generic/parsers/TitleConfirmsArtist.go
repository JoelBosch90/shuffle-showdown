package parsers

import "regexp"

func TitleConfirmsArtist(title string, artistName string) bool {
	sanitizedArtistName := regexp.QuoteMeta(artistName)
	regex := regexp.MustCompile("(?i).*" + sanitizedArtistName + ".*")
	match := regex.FindStringSubmatch(title)

	return match != nil
}
