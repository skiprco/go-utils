package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanSpecialCharacters(t *testing.T) {
	var tests = []struct {
		name   string
		input  string
		output string
	}{
		{"Remove special characters", "1111.222-B!?*5", "1111222B5"},
		{"Remove spaces", "1111		AA 2", "1111AA2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := CleanSpecialCharacters(tt.input)
			assert.Equal(t, tt.output, response)
		})
	}

}
