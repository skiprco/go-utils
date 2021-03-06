package validation

import (
	"github.com/nyaruka/phonenumbers"
	"github.com/skiprco/go-utils/v2/errors"
)

const phoneNumberInvalidCountryCodeMessage = "invalid country code"

// ValidateAndFormatPhoneNumber checks if the provided phone number is valid.
// If yes, it will format it to its E164 representation (e.g +32...).
//
// Raises
//
// - 400/invalid_country_code: The provided country code is invalid
//
// - 400/not_a_phone_number: The provided phone number is not recognised as one
//
// - 400/invalid_phone_number: The provided phone number has the correct format, but is symantically incorrect
//
// - 400/not_a_mobile_phone_number: The provided phone number is not a mobile number
func ValidateAndFormatPhoneNumber(phoneNumber string, countryCode string) (string, *errors.GenericError) {

	// This function will format any phoneNumber to its E164 representation (e.g +32...)
	//If the phoneNumber is already well formatted the parsing will succeed even with no country code
	parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber, countryCode)

	// If the parsing fails then it either means that the country code is required or the number is not valid at all
	if err != nil {
		if err.Error() == phoneNumberInvalidCountryCodeMessage {
			return "", errors.NewGenericError(400, errorDomain, errorSubDomain, ErrorInvalidCountryCode, nil)
		}
		return "", errors.NewGenericError(400, errorDomain, errorSubDomain, ErrorNotAPhoneNumber, nil)
	}
	if !phonenumbers.IsValidNumber(parsedPhoneNumber) {
		return "", errors.NewGenericError(400, errorDomain, errorSubDomain, ErrorInvalidPhoneNumber, nil)
	}
	phoneType := phonenumbers.GetNumberType(parsedPhoneNumber)
	if phoneType == phonenumbers.FIXED_LINE {
		return "", errors.NewGenericError(400, errorDomain, errorSubDomain, ErrorNotAMobilePhoneNumber, nil)
	}
	return phonenumbers.Format(parsedPhoneNumber, phonenumbers.E164), nil
}
