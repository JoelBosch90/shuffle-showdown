package formatters

func FormatTrackTitleWithFirstrArtistNameSongSuffix(trackTitle string, artistNames []string) string {
	cleanTitle := FormatCleanTrackTitle(trackTitle, artistNames)

	if len(artistNames) == 0 {
		return cleanTitle + " (song)"
	}

	return cleanTitle + " (" + artistNames[0] + " song)"
}
