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

func Test_SanitizeObject_Struct(t *testing.T) {
	type Test struct {
		TestPub  string
		testPriv string
		TestNil  *string
	}

	input := Test{
		TestPub:  "<p>Test</p>",
		testPriv: "<p>Test</p>",
		TestNil:  nil,
	}
	expected := Test{
		TestPub:  "Test",
		testPriv: "<p>Test</p>", // Private fields are unreachable
		TestNil:  nil,
	}
	genErr := SanitizeObject(&input)
	require.Nil(t, genErr)
	assert.Equal(t, expected, input)
}

func Test_SanitizeObject_Slice(t *testing.T) {
	input := []string{"<p>Test</p>"}
	expected := []string{"Test"}
	genErr := SanitizeObject(&input)
	require.Nil(t, genErr)
	assert.Equal(t, expected, input)
}

func Test_SanitizeObject_Map(t *testing.T) {
	// Should only sanitize map values. Keys should not be touched.
	input := map[string]string{"<p>Test</p>": "<p>Test</p>"}
	expected := map[string]string{"Test": "Test"}
	genErr := SanitizeObject(&input)
	require.Nil(t, genErr)
	assert.Equal(t, expected, input)
}
