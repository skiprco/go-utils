package metadata

import (
	"context"
	"testing"

	"github.com/micro/go-micro/v2/metadata"
	"github.com/skiprco/go-utils/v2/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetGoMicroMetadata_Empty(t *testing.T) {
	// Call helper
	result, genErr := GetGoMicroMetadata(context.Background())

	// Assert results
	require.Nil(t, genErr)
	assert.Equal(t, Metadata{}, result)
}

func Test_GetGoMicroMetadata_WithMetadata(t *testing.T) {
	// Setup test data
	ctx, _, genErr := UpdateGoMicroMetadata(context.Background(), fixtureMetadata())
	require.Nil(t, genErr)

	// Call helper
	result, genErr := GetGoMicroMetadata(ctx)

	// Assert results
	require.Nil(t, genErr)
	assert.Equal(t, fixtureMetadata(), result)
}

func Test_GetGoMicroMetadata_InvalidMetadata(t *testing.T) {
	// Call helper
	ctx := metadata.Set(context.Background(), GoMicroMetadataKey, "invalid")
	result, genErr := GetGoMicroMetadata(ctx)

	// Assert results
	assert.Nil(t, result)
	errors.AssertGenericError(t, genErr, 400, "decode_glob_from_base64_failed", nil)
}
