package gin

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/assert"
)

func Test_GetContextWithMetadata_WithMetadata(t *testing.T) {
	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}
	ginCtx.Set(contextMetadataKey, fixtureMetadata())

	// Call helper
	ctx := GetContextWithMetadata(ginCtx)

	// Assert result
	md, ok := metadata.FromContext(ctx)
	expected := metadata.Metadata{
		"Test_key_1": "test_key_1_success",
		"Test-Key-2": "test-key-2-success",
		"TestKey3":   "TestKey3Success",
	}
	assert.Equal(t, expected, md)
	assert.True(t, ok)
}

func Test_GetContextWithMetadata_WithoutMetadata(t *testing.T) {
	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}

	// Call helper
	ctx := GetContextWithMetadata(ginCtx)

	// Assert result
	md, ok := metadata.FromContext(ctx)
	assert.Nil(t, md)
	assert.False(t, ok)
}
