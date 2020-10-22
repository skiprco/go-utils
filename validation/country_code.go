package validation

import "github.com/skiprco/go-utils/converters"

// ValidateCountryCode checks if a country code is valid.
// An empty code is considered valid as well.
func ValidateCountryCode(code string) bool {
	// Empty code is valid
	if code == "" {
		return true
	}

	// Validate code
	_, genErr := converters.CountryCodeToCountryName(code)
	return genErr == nil
}
