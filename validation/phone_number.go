package validation

import (
	"github.com/nyaruka/phonenumbers"
	"github.com/skiprco/go-utils/v2/errors"
)

// ValidateAndFormatPhoneNumber checks if the provided phone number is valid.
// If yes, it will format it to its E164 representation (e.g +32...)
func ValidateAndFormatPhoneNumber(phoneNumber string, countryCode string) (string, *errors.GenericError) {
	// This function will format any phoneNumber to its E164 representation (e.g +32...)
	//If the phoneNumber is already well formatted the parsing will succeed even with no country code
	parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber, countryCode)

	// If the parsing fails then it either means that the country code is required or the number is not valid at all
	if err != nil {
		if countryCode == "" {
			return "", errors.NewGenericError(400, "go_utils", "common", "phone_number_country_code_required", nil)
		}
		return "", errors.NewGenericError(400, "go_utils", "common", "not_a_phone_number", nil)
	}
	if !phonenumbers.IsValidNumber(parsedPhoneNumber) {
		return "", errors.NewGenericError(400, "go_utils", "common", "invalid_phone_number", nil)
	}
	phoneType := phonenumbers.GetNumberType(parsedPhoneNumber)
	if phoneType == phonenumbers.FIXED_LINE {
		return "", errors.NewGenericError(400, "go_utils", "common", "not_a_mobile_phone_number", nil)
	}
	return phonenumbers.Format(parsedPhoneNumber, phonenumbers.E164), nil
}
