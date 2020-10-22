package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateCountryCode_Success(t *testing.T) {
	result := ValidateCountryCode("")
	assert.True(t, result)
}

func Test_ValidateCountryCode_Empty_Success(t *testing.T) {
	result := ValidateCountryCode("BE")
	assert.True(t, result)
}

func Test_ValidateCountryCode_Failure(t *testing.T) {
	result := ValidateCountryCode("XX")
	assert.False(t, result)
}
