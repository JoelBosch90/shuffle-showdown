package helpers

func IncludesString(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}

	return false
}

func IncludesStrings(haystack []string, needles []string) bool {
	for _, needle := range needles {
		if !IncludesString(haystack, needle) {
			return false
		}
	}

	return true
}
