package converters

import (
	"testing"

	"github.com/skiprco/go-utils/v2/errors"
	"github.com/stretchr/testify/assert"
)

func Test_CountryCodeToCountryName_Success(t *testing.T) {
	result, genErr := CountryCodeToCountryName("BE")
	assert.Nil(t, genErr)
	assert.Equal(t, "Belgium", result)
}

func Test_CountryCodeToCountryName_Failure(t *testing.T) {
	result, genErr := CountryCodeToCountryName("XX")
	assert.Equal(t, "", result)
	errors.AssertGenericError(t, genErr, 404, "country_not_found", nil)
}

func Test_CountryNameToCountryCode_Simple_Success(t *testing.T) {
	result, genErr := CountryNameToCountryCode("Belgium")
	assert.Nil(t, genErr)
	assert.Equal(t, "BE", result)
}

func Test_CountryNameToCountryCode_Complex_Success(t *testing.T) {
	result, genErr := CountryNameToCountryCode("CURAÇAO")
	assert.Nil(t, genErr)
	assert.Equal(t, "CW", result)
}

func Test_CountryNameToCountryCode_Failure(t *testing.T) {
	result, genErr := CountryNameToCountryCode("Invalid")
	assert.Equal(t, "", result)
	errors.AssertGenericError(t, genErr, 404, "country_not_found", nil)
}

func Test_IsValidCountryCode_Success(t *testing.T) {
	result := IsValidCountryCode("BE")
	assert.True(t, result)
}

func Test_IsValidCountryCode_Failure(t *testing.T) {
	result := IsValidCountryCode("BEL")
	assert.False(t, result)
}
