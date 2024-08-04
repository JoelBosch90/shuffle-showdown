package formatters

func FormatTrackTitleWithOnlySongSuffix(trackTitle string, _ []string) string {
	return FormatTrackTitleWithArtistNamesSongSuffix(trackTitle, []string{})
}
