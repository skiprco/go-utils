package converters

import (
	"github.com/skiprco/go-utils/v2/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ConvertCentimeWithEnLanguageAndUSD_success(t *testing.T) {
	result, genErr := ConvertCentToCurrency(200, "USD", "EN")
	assert.Nil(t, genErr)
	assert.Equal(t, "$2.00", result)
}

func Test_ConvertCentimeWithFRLanguageAndEUR_success(t *testing.T) {
	result, genErr := ConvertCentToCurrency(200, "EUR", "FR")
	assert.Nil(t, genErr)
	assert.Equal(t, "2,00 €", result)
}

func Test_ConvertCentimeWithUnknownLanguageAndEUR_success(t *testing.T) {
	result, genErr := ConvertCentToCurrency(200, "EUR", "")
	assert.Nil(t, genErr)
	assert.Equal(t, "2.00 €", result)
}

func Test_ConvertCentimeWithUnknownCurrency_failed500(t *testing.T) {
	result, genErr := ConvertCentToCurrency(200, "OOO", "FR")
	assert.NotNil(t, errors.NewGenericError(400, errorDomain, errorSubDomain, ErrorUnknownCurrency, map[string]string{"currency": "OOO"}), genErr)
	assert.Equal(t, "", result)
}
