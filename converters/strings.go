package converters

import (
	"regexp"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/v2/errors"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormaliseString removes all the accents from the letters in the string.
// Based on https://twinnation.org/articles/33/remove-accents-from-characters-in-go
//
// Raises
//
// - 400/failed_to_normalise_string: Provided country name could not be normalised
func NormaliseString(input string) (string, *errors.GenericError) {
	// Normalise string
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, input)

	// Handle error
	if err != nil {
		log.WithField("input", input).WithField("error", err).Error("Unable to normalise string")
		return "", errors.NewGenericError(400, errorDomain, errorSubDomain, ErrorFailedToNormaliseString, nil)
	}

	// Normalise successful
	return output, nil
}

var matchFirstCap = regexp.MustCompile("([A-Z])([A-Z][a-z])")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts the provided string to snake_case.
// Based on https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
func ToSnakeCase(input string) string {
	output := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	output = matchAllCap.ReplaceAllString(output, "${1}_${2}")
	output = strings.ReplaceAll(output, "-", "_")
	return strings.ToLower(output)
}

// CleanSpecialCharacters removes any character(included spaces) which is
// not a digit or a letter from the input.
//
// Raises
//
// Nothing: This function will never raise an error
func CleanSpecialCharacters(input string) (string, *errors.GenericError) {
	normalisedInput, genErr := NormaliseString(input)
	if genErr != nil {
		return "", genErr
	}
	regExp := regexp.MustCompile(`[^0-9a-zA-Z]+`)
	return regExp.ReplaceAllLiteralString(normalisedInput, ""), nil
}
