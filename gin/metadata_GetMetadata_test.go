package gin

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_GetMetadata_Empty(t *testing.T) {
	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}

	// Call helper
	result := GetMetadata(ginCtx)

	// Assert results
	assert.Equal(t, map[string]string{}, result)
}

func Test_GetMetadata_WithMetadata(t *testing.T) {
	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}
	ginCtx.Set(contextMetadataKey, fixtureMetadata())

	// Call helper
	result := GetMetadata(ginCtx)

	// Assert results
	assert.Equal(t, fixtureMetadata(), result)
}
