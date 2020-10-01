package errors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertGenericError helps to assert GenericErrors during unit testing
func AssertGenericError(t *testing.T, err error, code int, detailContains string, metaContains map[string]string) {
	// Assert type
	if assert.NotNil(t, err, "GenericError should not be nil") && assert.IsType(t, &GenericError{}, err) {
		// Assert base fields
		genericError := err.(*GenericError)
		assert.Equal(t, code, genericError.Code)
		assert.Contains(t, strings.ToLower(genericError.GetDetailString()), strings.ToLower(detailContains))

		// Assert meta
		if metaContains != nil {
			for metaKey, metaValue := range metaContains {
				if assert.Contains(t, genericError.Meta, metaKey) {
					assert.Equal(t, metaValue, genericError.Meta[metaKey])
				}
			}
		}
	}
}
