package validation

import (
	"time"
)

// will be move and tested in go-utils package
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

//  checks if a time is within a given timerange
// an empty end means infinity
func timeWithinTimerage(now time.Time, start time.Time, end time.Time) bool {
	// Empty code is valid
	if code == "" {
		return true
	}

	// Validate code
	_, genErr := converters.CountryCodeToCountryName(code)
	return genErr == nil
}
