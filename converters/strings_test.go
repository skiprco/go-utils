package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NormaliseString_EmptyInput(t *testing.T) {
	result, genErr := NormaliseString("")
	assert.Nil(t, genErr)
	assert.Equal(t, "", result)
}

func Test_NormaliseString_NormalInput(t *testing.T) {
	result, genErr := NormaliseString("Tëst Çôdé (should-also-work)")
	assert.Nil(t, genErr)
	assert.Equal(t, "Test Code (should-also-work)", result)
}
