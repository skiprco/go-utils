package metadata

import (
	"context"
	"testing"

	"github.com/micro/go-micro/v2/metadata"
	"github.com/skiprco/go-utils/v2/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UpdateGoMicroMetadata_NoMetadataYet(t *testing.T) {
	// Should only return additional metadata

	// Call helper
	ctx, meta, genErr := UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Assert results
	require.Nil(t, genErr)
	assert.Equal(t, fixtureMetadata(), meta)
	result, genErr := GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, fixtureMetadata(), result)
}

func Test_UpdateGoMicroMetadata_HasAlreadyMetadata(t *testing.T) {
	// Should merge current and additional metadata

	// Setup test data
	ctx, _, _ := UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Call helper
	additionalMeta := Metadata{
		"PascalKey": "PascalValueNew",
		"NewKey":    "new_value",
	}
	ctx, meta, genErr := UpdateGoMicroMetadata(ctx, additionalMeta)

	// Assert results
	expected := fixtureMetadata()
	expected["PascalKey"] = "PascalValueNew"
	expected["NewKey"] = "new_value"
	require.Nil(t, genErr)
	assert.Equal(t, expected, meta)
	result, genErr := GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, expected, result)
}

func Test_UpdateGoMicroMetadata_AdditionalMetadataNil(t *testing.T) {
	// Should keep original set metadata

	// Setup test data
	ctx, _, _ := UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Call helper
	ctx, meta, genErr := UpdateGoMicroMetadata(ctx, nil)

	// Assert results
	expected := fixtureMetadata()
	require.Nil(t, genErr)
	assert.Equal(t, expected, meta)
	result, genErr := GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, expected, result)
}

func Test_UpdateGoMicroMetadata_InvalidMetadata(t *testing.T) {
	// Call helper
	reqCtx := metadata.Set(context.Background(), GoMicroMetadataKey, "invalid")
	resCtx, meta, genErr := UpdateGoMicroMetadata(reqCtx, fixtureMetadata())

	// Assert results
	assert.Same(t, reqCtx, resCtx)
	assert.Nil(t, meta)
	errors.AssertGenericError(t, genErr, 400, "decode_glob_from_base64_failed", nil)
}
