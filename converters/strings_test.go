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

func Test_ToSnakeCase(t *testing.T) {
	// Based on https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"camelCase", "camel_case"},
		{"PascalCase", "pascal_case"},
		{"snake_case", "snake_case"},
		{"Pascal_Snake", "pascal_snake"},
		{"SCREAMING_SNAKE", "screaming_snake"},
		{"kebab-case", "kebab_case"},
		{"Pascal-Kebab", "pascal_kebab"},
		{"SCREAMING-KEBAB", "screaming_kebab"},
		{"A", "a"},
		{"AA", "aa"},
		{"AAA", "aaa"},
		{"AAAA", "aaaa"},
		{"AaAa", "aa_aa"},
		{"HTTPRequest", "http_request"},
		{"BatteryLifeValue", "battery_life_value"},
		{"Id0Value", "id0_value"},
		{"ID0Value", "id0_value"},
	}

	for _, test := range tests {
		result := ToSnakeCase(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func Test_CleanSpecialCharacters(t *testing.T) {
	var tests = []struct {
		name   string
		input  string
		output string
	}{
		{"Input already clean", "1111AA2", "1111AA2"},
		{"Remove special characters", "1111.222-B!?*5", "1111222B5"},
		{"Remove spaces", "1111		AA 2", "1111AA2"},
		{"Normalise input", "1111AA 2Ç", "1111AA2C"},
		{"Empty input", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := CleanSpecialCharacters(tt.input)
			assert.Equal(t, tt.output, response)
		})
	}

}
