package formatters

func FormatTrackTitleWithSongSuffix(trackTitle string, artistNames []string) string {
	return FormatCleanTrackTitle(trackTitle, artistNames) + " (lied)"
}
