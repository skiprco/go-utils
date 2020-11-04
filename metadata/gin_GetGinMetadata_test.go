package metadata

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_GetGinMetadata_Empty(t *testing.T) {
	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}

	// Call helper
	result := GetGinMetadata(ginCtx)

	// Assert results
	assert.Equal(t, Metadata{}, result)
}

func Test_GetGinMetadata_WithMetadata(t *testing.T) {
	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}
	UpdateGinMetadata(ginCtx, fixtureMetadata())

	// Call helper
	result := GetGinMetadata(ginCtx)

	// Assert results
	assert.Equal(t, fixtureMetadata(), result)
}
