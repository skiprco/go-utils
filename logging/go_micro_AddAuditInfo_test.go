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

func Test_AddAuditInfo_New(t *testing.T) {
	// Setup test data
	ctx, _, _ := metadata.SetGoMicroMetadata(context.Background(), "service_name", "srv-test")

	// Call helper
	ctx, genErr := AddAuditInfo(ctx, "test_key", "test_value")

	// Assert results
	expected := metadata.Metadata{
		"service_name":      "srv-test",
		"srv_test_test_key": "test_value",
	}
	require.Nil(t, genErr)
	meta, genErr := metadata.GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, expected, meta)
}

func Test_AddAuditInfo_Overwrite(t *testing.T) {
	// Setup test data
	ctx, _, _ := metadata.UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Call helper
	ctx, genErr := AddAuditInfo(ctx, "test_key", "after_value")

	// Assert results
	expected := fixtureMetadata()
	expected["srv_test_test_key"] = "after_value"
	require.Nil(t, genErr)
	meta, genErr := metadata.GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, expected, meta)
}

func Test_AddAuditInfo_InvalidMeta(t *testing.T) {
	// Setup test data
	ctx := microMeta.Set(context.Background(), metadata.GoMicroMetadataKey, "invalid")

	// Call helper
	ctx, genErr := AddAuditInfo(ctx, "test_key", "test_value")

	// Assert results
	assert.NotNil(t, genErr)
	assert.NotNil(t, ctx)
}

func Test_AddAuditInfo_ServiceNameNotSet(t *testing.T) {
	// Call helper
	ctx, genErr := AddAuditInfo(context.Background(), "test_key", "test_value")

	// Assert results
	assert.NotNil(t, ctx)
	meta := map[string]string{"key": "service_name"}
	errors.AssertGenericError(t, genErr, 500, "key_not_found_in_context", meta)
}
