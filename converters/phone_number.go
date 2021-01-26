package converters

import (
	"github.com/skiprco/go-utils/v2/errors"
	"regexp"
)

/**
CleanPhoneNumber remove all non necessary code of a phone number
Be aware that remove all non numeric char except the sign '+'
*/
func CleanPhoneNumber(phoneNumber string) (string, *errors.GenericError) {
	regExpression := regexp.MustCompile(`[^\+\d]`)
	return regExpression.ReplaceAllLiteralString(phoneNumber, ""), nil
}
