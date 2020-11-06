package logging

import (
	"context"
	"testing"

	"github.com/skiprco/go-utils/v2/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AddAuditInfo_New(t *testing.T) {
	ctx, genErr := AddAuditInfo(context.Background(), "Test", "test_key", "test_value")

	// Assert results
	expected := metadata.Metadata{"test_test_key": "test_value"}
	require.Nil(t, genErr)
	meta, genErr := metadata.GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, expected, meta)
}

func Test_AddAuditInfo_Overwrite(t *testing.T) {
	// Setup test data
	ctx, _, _ := metadata.UpdateGoMicroMetadata(context.Background(), fixtureMetadata())

	// Call helper
	ctx, genErr := AddAuditInfo(ctx, "TestService", "test_key", "after_value")

	// Assert results
	require.Nil(t, genErr)
	meta, genErr := metadata.GetGoMicroMetadata(ctx)
	require.Nil(t, genErr)
	assert.Equal(t, fixtureMetadata(), meta)
}
