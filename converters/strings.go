package converters

import (
	"unicode"

	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/errors"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormaliseString removes all the accents from the letters in the string.
// Based on https://twinnation.org/articles/33/remove-accents-from-characters-in-go
func NormaliseString(input string) (string, *errors.GenericError) {
	// Normalise string
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, input)

	// Handle error
	if err != nil {
		log.WithField("input", input).WithField("error", err).Error("Unable to normalise string")
	}

	// Normalise successful
	return output, nil
}
