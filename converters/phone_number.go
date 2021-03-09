package converters

import (
	"regexp"

	"github.com/skiprco/go-utils/v3/errors"
)

// CleanPhoneNumber remove all non necessary code of a phone number
// Be aware that remove all non numeric char except the sign '+'
//
// Raises
//
// Nothing: This function will never raise an error
func CleanPhoneNumber(phoneNumber string) (string, *errors.GenericError) {
	regExpression := regexp.MustCompile(`[^\+\d]`)
	return regExpression.ReplaceAllLiteralString(phoneNumber, ""), nil
}
