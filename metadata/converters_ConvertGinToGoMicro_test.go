package metadata

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ConvertGinToGoMicro_WithMetadata(t *testing.T) {
	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}
	UpdateGinMetadata(ginCtx, fixtureMetadata())

	// Call helper
	ctx := ConvertGinToGoMicro(ginCtx)

	// Assert result
	meta, genErr := GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, fixtureMetadata(), meta)
}

func Test_ConvertGinToGoMicro_WithoutMetadata(t *testing.T) {
	// Setup test data
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}

	// Call helper
	ctx := ConvertGinToGoMicro(ginCtx)

	// Assert result
	meta, genErr := GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, Metadata{}, meta)
}
