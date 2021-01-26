package converters

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CleanPhoneNumber_success(t *testing.T) {
	result, err := CleanPhoneNumber("+32475000001")
	assert.Equal(t, "+32475000001", result)
	require.Nil(t, err)
}

func Test_CleanPhoneNumberWithInvisibleChar_success(t *testing.T) {
	result, err := CleanPhoneNumber("+32475000001​")
	assert.Equal(t, "+32475000001", result)
	require.Nil(t, err)
}

func Test_CleanPhoneNumberWithUnicode_success(t *testing.T) {
	result, err := CleanPhoneNumber("+32475000001\u200b")
	assert.Equal(t, "+32475000001", result)
	require.Nil(t, err)
}

func Test_CleanPhoneNumberWithEmpty_success(t *testing.T) {
	result, err := CleanPhoneNumber("")
	assert.Equal(t, "", result)
	require.Nil(t, err)
}
