package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPhoneNumber(t *testing.T) {
	var tests = []struct {
		name         string
		input        string
		inputCountry string
		output       string
	}{
		{"well formatted phone number BE", "+32468300431", "BE", "+32468300431"},
		{"well formatted phone number BE without passing country code", "+32468300431", "", "+32468300431"},
		{"Valid phone number FR", "06 23 83 96 79", "FR", "+33623839679"},
		{"Valid phone number but not correctly formatted", "0468300431", "BE", "+32468300431"},
		{"Valid phone number but not correctly formatted without passing country code", "0468300431", "", ""},
		{"Not Valid phone number", "0461", "BE", ""},
		{"Valid phone number but not a mobile", "+3227896143", "BE", ""},
		{"Not a phone number at all", "dhjdhj", "BE", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := ValidateAndFormatPhoneNumber(tt.input, tt.inputCountry)
			assert.Equal(t, tt.output, response)
		})
	}

}
