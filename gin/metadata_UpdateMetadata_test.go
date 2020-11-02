package gin

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UpdateMetadata_NoMetadataYet(t *testing.T) {
	// Should only return additional metadata

	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}

	// Call helper
	result := UpdateMetadata(ginCtx, fixtureMetadata())

	// Assert results
	require.Equal(t, fixtureMetadata(), result)
	assert.Equal(t, GetMetadata(ginCtx), result)
}

func Test_UpdateMetadata_HasAlreadyMetadata(t *testing.T) {
	// Should merge current and additional metadata

	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}
	ginCtx.Set(contextMetadataKey, fixtureMetadata())

	// Call helper
	additionalMeta := map[string]string{
		"test_key_1": "test_key_1_success_new",
		"test_key_4": "test_key_4_success",
	}
	result := UpdateMetadata(ginCtx, additionalMeta)

	// Assert results
	expected := fixtureMetadata()
	expected["test_key_1"] = "test_key_1_success_new"
	expected["test_key_4"] = "test_key_4_success"
	require.Equal(t, expected, result)
	assert.Equal(t, GetMetadata(ginCtx), result)
}

func Test_UpdateMetadata_AdditionalMetadataNil(t *testing.T) {
	// Should keep original set metadata

	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}
	ginCtx.Set(contextMetadataKey, fixtureMetadata())

	// Call helper
	result := UpdateMetadata(ginCtx, nil)

	// Assert results
	require.Equal(t, fixtureMetadata(), result)
	assert.Equal(t, GetMetadata(ginCtx), result)
}
