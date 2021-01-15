package converters

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CleanPhoneNumber_success(t *testing.T){
	result := CleanPhoneNumber("+32475000001")
	assert.Equal(t, "+32475000001", result)
}

func Test_CleanPhoneNumberWithInvisibleChar_success(t *testing.T){
	result := CleanPhoneNumber("+32475000001​")
	assert.Equal(t, "+32475000001", result)
}

func Test_CleanPhoneNumberWithUnicode_success(t *testing.T){
	result := CleanPhoneNumber("+32475000001\u200b")
	assert.Equal(t, "+32475000001", result)
}

func Test_CleanPhoneNumberWithEmpty_success(t *testing.T){
	result := CleanPhoneNumber("")
	assert.Equal(t, "", result)
}
