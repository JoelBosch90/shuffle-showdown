package upsert

import (
	"strings"
)

func isCapitalLetter(character rune) bool {
	return 'A' <= character && character <= 'Z'
}

func PascalToSnake(pascalText string) string {
	var snakeText string

	for index, currentCharacter := range pascalText {
		// If the first character is a capital letter, make it lowercase.
		if index == 0 {
			snakeText += strings.ToLower(string(currentCharacter))
			continue
		}

		// If the previous character is a lower case letter, and the
		// current character is an upper case letter, add an underscore.
		previousCharacter := []rune(pascalText)[index-1]
		if !isCapitalLetter(previousCharacter) && isCapitalLetter(currentCharacter) {
			snakeText += "_"
		}

		snakeText += strings.ToLower(string(currentCharacter))
	}

	return snakeText
}
