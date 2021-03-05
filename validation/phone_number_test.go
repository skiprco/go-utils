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
		err          string
	}{
		{"well formatted phone number BE", "+32468300431", "BE", "+32468300431", ""},
		{"well formatted phone number BE without passing country code", "+32468300431", "", "+32468300431", ""},
		{"Valid phone number FR", "06 23 83 96 79", "FR", "+33623839679", ""},
		{"Valid phone number but not correctly formatted", "0468300431", "BE", "+32468300431", ""},
		{"Valid phone number but not correctly formatted without passing country code", "0468300431", "", "", ErrorInvalidCountryCode},
		{"Not Valid phone number", "0461", "BE", "", ErrorInvalidPhoneNumber},
		{"Valid phone number but not a mobile", "+3227896143", "BE", "", ErrorNotAMobilePhoneNumber},
		{"Not a phone number at all", "dhjdhj", "BE", "", ErrorNotAPhoneNumber},
		{"Not a phone number at all without country code", "+324ER3643", "", "", ErrorInvalidPhoneNumber},
		{"Invalid phone number without country code", "+324780000012", "", "", ErrorInvalidPhoneNumber},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := ValidateAndFormatPhoneNumber(tt.input, tt.inputCountry)
			assert.Equal(t, tt.output, response)
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err, err.SubDomainCode)
			}
		})
	}

}
