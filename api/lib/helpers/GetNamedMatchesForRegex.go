package helpers

import (
	"regexp"
)

func GetNamedMatchesForRegex(regex *regexp.Regexp, name string) map[string]string {
	matches := make(map[string]string)
	match := regex.FindStringSubmatch(name)

	if match == nil {
		return matches
	}

	for index, name := range regex.SubexpNames() {
		if index != 0 && name != "" {
			matches[name] = match[index]
		}
	}

	return matches
}
