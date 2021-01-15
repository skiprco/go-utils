package converters

import (
	"regexp"
)

/**
CleanPhoneNumber remove all non necessary code of a phone number
Be aware that remove all non numeric char except the sign '+'
 */
func CleanPhoneNumber(phoneNumber string) string {
	regExpression := regexp.MustCompile(`[^\+\d]`)
	return regExpression.ReplaceAllLiteralString(phoneNumber, "")
}
