package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sanitize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"Test", "Test"},
		{"Skipr&Co", "Skipr&Co"},
		{"O'Hare", "O'Hare"},
		{"<img src='test' />", ""},
		{"<p>Test</p>", "Test"},
	}

	for _, test := range tests {
		result := Sanitize(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func Test_SanitizeObject_String(t *testing.T) {
	input := "<p>Test</p>"
	expected := "Test"
	genErr := SanitizeObject(&input)
	require.Nil(t, genErr)
	assert.Equal(t, expected, input)
}

func Test_SanitizeObject_StringPointer(t *testing.T) {
	input := "<p>Test</p>"
	inputPtr := &input
	expected := "Test"
	expectedPtr := &expected
	genErr := SanitizeObject(&inputPtr)
	require.Nil(t, genErr)
	assert.Equal(t, expectedPtr, inputPtr)
}
