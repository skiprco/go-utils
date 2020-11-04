package logging

import (
	"context"
	"testing"

	"github.com/skiprco/go-utils/v2/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AddAuditInfoMap_WithData(t *testing.T) {
	// Setup test data
	ctx, _, _ := metadata.UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Call helper
	ctx, genErr := AddAuditInfoMap(ctx, "TestService", metadata.Metadata{
		"TestKey": "after_value",
		"NewKey":  "new_value",
	})

	// Assert results
	expected := fixtureMetadata()
	expected["test_service_test_key"] = "after_value"
	expected["test_service_new_key"] = "new_value"
	require.Nil(t, genErr)
	meta, genErr := metadata.GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, expected, meta)
}

func Test_AddAuditInfoMap_Nil(t *testing.T) {
	// Setup test data
	ctx, _, _ := metadata.UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Call helper
	ctx, genErr := AddAuditInfoMap(ctx, "TestService", nil)

	// Assert results
	require.Nil(t, genErr)
	meta, genErr := metadata.GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, fixtureMetadata(), meta)
}
