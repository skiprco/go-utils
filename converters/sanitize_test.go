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

func Test_SanitizeObject_Complex(t *testing.T) {
	// Define test specific types and helpers
	type Item struct {
		TestString string
		TestInt    int
		testPriv   string
	}

	type Nested struct {
		Test          string
		ItemPtr       *Item
		StringItemMap map[string]Item
	}

	type Base struct {
		Nested
		ItemItemMap map[Item]*Item
		ItemSlice   []*Item
	}

	dirtyItem := func() *Item {
		return &Item{
			TestString: "<p>Test</p>",
			TestInt:    8,
			testPriv:   "<p>Test</p>",
		}
	}

	cleanItem := func() *Item {
		return &Item{
			TestString: "Test",
			TestInt:    8,
			testPriv:   "<p>Test</p>", // Private fields are unreachable
		}
	}

	// Define test data
	input := Base{
		ItemItemMap: map[Item]*Item{
			*dirtyItem(): dirtyItem(),
		},
		ItemSlice: []*Item{
			dirtyItem(),
			dirtyItem(),
			nil,
		},
	}
	input.Nested = Nested{
		Test:    "<p>Test</p>",
		ItemPtr: dirtyItem(),
		StringItemMap: map[string]Item{
			"<p>Test</p>": *dirtyItem(),
			"Test2":       *dirtyItem(),
		},
	}

	// Define expected
	expected := Base{
		ItemItemMap: map[Item]*Item{
			*cleanItem(): cleanItem(),
		},
		ItemSlice: []*Item{
			cleanItem(),
			cleanItem(),
			nil,
		},
	}
	expected.Nested = Nested{
		Test:    "Test",
		ItemPtr: cleanItem(),
		StringItemMap: map[string]Item{
			"Test":  *cleanItem(),
			"Test2": *cleanItem(),
		},
	}

	// Call helper and validate result
	genErr := SanitizeObject(&input)
	require.Nil(t, genErr)
	assert.Equal(t, expected, input)
}
