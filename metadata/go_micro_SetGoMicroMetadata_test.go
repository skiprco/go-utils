package metadata

import (
	"context"
	"testing"

	"github.com/asim/go-micro/v3/metadata"
	"github.com/skiprco/go-utils/v3/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SetGoMicroMetadata_NoMetadataYet(t *testing.T) {
	// Should only return additional metadata

	// Call helper
	ctx, meta, genErr := SetGoMicroMetadata(context.Background(), "test_key", "test_value")

	// Assert results
	expected := Metadata{"test_key": "test_value"}
	require.Nil(t, genErr)
	assert.Equal(t, expected, meta)
	result, genErr := GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, expected, result)
}

func Test_SetGoMicroMetadata_HasAlreadyMetadata(t *testing.T) {
	// Should merge current and additional metadata

	// Setup test data
	ctx, _, _ := UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Call helper
	ctx, meta, genErr := SetGoMicroMetadata(ctx, "PascalKey", "PascalValueNew")

	// Assert results
	expected := fixtureMetadata()
	expected["PascalKey"] = "PascalValueNew"
	require.Nil(t, genErr)
	assert.Equal(t, expected, meta)
	result, genErr := GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, expected, result)
}

func Test_SetGoMicroMetadata_InvalidMetadata(t *testing.T) {
	// Call helper
	reqCtx := metadata.Set(context.Background(), GoMicroMetadataKey, "invalid")
	resCtx, meta, genErr := SetGoMicroMetadata(reqCtx, "test_key", "test_value")

	// Assert results
	assert.Same(t, reqCtx, resCtx)
	assert.Nil(t, meta)
	errors.AssertGenericError(t, genErr, 400, "decode_glob_from_base64_failed", nil)
}
