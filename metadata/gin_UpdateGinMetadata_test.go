package metadata

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateGinMetadata_NoMetadataYet(t *testing.T) {
	// Should only return additional metadata

	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}

	// Call helper
	result := UpdateGinMetadata(ginCtx, fixtureMetadata())

	// Assert results
	expected := fixtureMetadata()
	assert.Equal(t, expected, result)
	assert.Equal(t, expected, GetGinMetadata(ginCtx))
}

func Test_UpdateGinMetadata_HasAlreadyMetadata(t *testing.T) {
	// Should merge current and additional metadata

	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}
	UpdateGinMetadata(ginCtx, fixtureMetadata())

	// Call helper
	additionalMeta := Metadata{
		"PascalKey": "PascalValueNew",
		"NewKey":    "new_value",
	}
	result := UpdateGinMetadata(ginCtx, additionalMeta)

	// Assert results
	expected := fixtureMetadata()
	expected["PascalKey"] = "PascalValueNew"
	expected["NewKey"] = "new_value"
	assert.Equal(t, expected, result)
	assert.Equal(t, expected, GetGinMetadata(ginCtx))
}

func Test_UpdateGinMetadata_AdditionalMetadataNil(t *testing.T) {
	// Should keep original set metadata

	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}
	UpdateGinMetadata(ginCtx, fixtureMetadata())

	// Call helper
	result := UpdateGinMetadata(ginCtx, nil)

	// Assert results
	expected := fixtureMetadata()
	assert.Equal(t, expected, result)
	assert.Equal(t, expected, GetGinMetadata(ginCtx))
}
