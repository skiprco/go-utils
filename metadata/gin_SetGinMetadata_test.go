package metadata

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_SetGinMetadata_NoMetadataYet(t *testing.T) {
	// Should only return additional metadata

	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}

	// Call helper
	result := SetGinMetadata(ginCtx, "test_key", "test_value")

	// Assert results
	expected := Metadata{"test_key": "test_value"}
	assert.Equal(t, expected, result)
	assert.Equal(t, expected, GetGinMetadata(ginCtx))
}

func Test_SetGinMetadata_HasAlreadyMetadata(t *testing.T) {
	// Should merge current and additional metadata

	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}
	UpdateGinMetadata(ginCtx, fixtureMetadata())

	// Call helper
	result := SetGinMetadata(ginCtx, "PascalKey", "PascalValueNew")

	// Assert results
	expected := fixtureMetadata()
	expected["PascalKey"] = "PascalValueNew"
	assert.Equal(t, expected, result)
	assert.Equal(t, expected, GetGinMetadata(ginCtx))
}
