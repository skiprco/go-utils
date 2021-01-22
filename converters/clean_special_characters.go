package converters

import (
	"regexp"
)

/*
CleanSpecialCharacters remove all special characters from a string input
*/
func CleanSpecialCharacters(input string) string {
	regExp := regexp.MustCompile(`[^0-9a-zA-Z]+`)
	return regExp.ReplaceAllLiteralString(input, "")
}
