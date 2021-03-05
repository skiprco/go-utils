package logging

import (
	"context"
	"testing"

	microMeta "github.com/micro/go-micro/v2/metadata"
	"github.com/skiprco/go-utils/v2/errors"
	"github.com/skiprco/go-utils/v2/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AddAuditInfoMap_WithData(t *testing.T) {
	// Setup test data
	ctx, _, _ := metadata.UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Call helper
	ctx, genErr := AddAuditInfoMap(ctx, metadata.Metadata{
		"TestKey": "after_value",
		"NewKey":  "new_value",
	})

	// Assert results
	expected := fixtureMetadata()
	expected["srv_test_test_key"] = "after_value"
	expected["srv_test_new_key"] = "new_value"
	require.Nil(t, genErr)
	meta, genErr := metadata.GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, expected, meta)
}

func Test_AddAuditInfoMap_Nil(t *testing.T) {
	// Setup test data
	ctx, _, _ := metadata.UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Call helper
	ctx, genErr := AddAuditInfoMap(ctx, nil)

	// Assert results
	require.Nil(t, genErr)
	meta, genErr := metadata.GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, fixtureMetadata(), meta)
}

func Test_AddAuditInfoMap_InvalidMeta(t *testing.T) {
	// Setup test data
	ctx := microMeta.Set(context.Background(), metadata.GoMicroMetadataKey, "invalid")

	// Call helper
	ctx, genErr := AddAuditInfoMap(ctx, nil)

	// Assert results
	assert.NotNil(t, genErr)
	assert.NotNil(t, ctx)
}

func Test_AddAuditInfoMap_ServiceNameNotSet(t *testing.T) {
	// Call helper
	ctx, genErr := AddAuditInfoMap(context.Background(), nil)

	// Assert results
	assert.NotNil(t, ctx)
	meta := map[string]string{"key": "service_name"}
	errors.AssertGenericError(t, genErr, 500, "key_not_found_in_context", meta)
}
