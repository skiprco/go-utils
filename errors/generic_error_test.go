package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericError_GetDetailStringWithoutMeta(t *testing.T) {
	// Get error string
	testError := NewGenericError(418, "test_domain", "test_subdomain", "test_error", nil)
	testError.Meta = map[string]string{}

	// Assert result
	assert.Equal(t, "test_domain/test_subdomain/test_error/", testError.GetDetailString())
}

func TestGenericError_GetDetailStringWithMeta(t *testing.T) {
	// Set defaults
	defaultMeta = map[string]string{
		"provider": "test_provider",
	}
	SetupDefaults(defaultMeta)

	// Get error string
	meta := map[string]string{"additional": "success"}
	detailString := NewGenericError(418, "test_domain", "test_subdomain", "test_error", meta).GetDetailString()

	// Assert result
	assert.Regexp(t, `^test_domain/test_subdomain/test_error/.+=.+(;.+?=.+?)*$`, detailString)
	assert.Contains(t, detailString, `provider=test_provider`)
	assert.Contains(t, detailString, `additional=success`)
}

func TestGenericError_ConvertToString(t *testing.T) {
	// Set defaults
	defaultMeta = map[string]string{
		"provider": "test_provider",
	}
	SetupDefaults(defaultMeta)

	// Get error string
	meta := map[string]string{"additional": "success"}
	errorString := NewGenericError(418, "test_domain", "test_subdomain", "test_error", meta).Error()

	// Assert result
	assert.Regexp(t, `"detail": test_domain/test_subdomain/test_error/.+=.+;.+=.+`, errorString)
	assert.Contains(t, errorString, `provider=test_provider`)
	assert.Contains(t, errorString, `additional=success`)

}
