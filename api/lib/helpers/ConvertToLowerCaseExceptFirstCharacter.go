package helpers

import "strings"

func ConvertToLowerCaseExceptFirstCharacter(artistName string) string {
	return strings.ToUpper(artistName[0:1]) + strings.ToLower(artistName[1:])
}
