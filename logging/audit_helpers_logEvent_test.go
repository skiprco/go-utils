package logging

import (
	"context"
	"testing"

	microMetadata "github.com/micro/go-micro/v2/metadata"
	log "github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/skiprco/go-utils/v2/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_logEvent_WithContextAndAdditionalData(t *testing.T) {
	// This test checks if the log data is correctly
	// merged from the different sources (context, parameters, ...)

	// Setup test
	hook := logTest.NewGlobal()
	meta := metadata.Metadata{
		"TestMetaKey":          "test-meta-value",
		"TestAdditionalString": "should-be-overwritten",
	}
	ctx, _, _ := metadata.UpdateGoMicroMetadata(context.Background(), meta)
	additional := map[string]interface{}{
		"test-additional-string": "test-additional-value",
		"test_additional_int":    123,
		"category":               "should-be-overwritten",
	}

	// Call helper
	logEvent(ctx, "test-message", AuditCategorySuccess, additional)

	// Assert result
	require.Len(t, hook.Entries, 1)
	expectedData := log.Fields{
		"category":               AuditCategorySuccess,
		"test_additional_string": "test-additional-value",
		"test_additional_int":    123,
		"test_meta_key":          "test-meta-value",
	}
	assert.Equal(t, expectedData, hook.LastEntry().Data)
	assert.Equal(t, "test-message", hook.LastEntry().Message)
	hook.Reset()
}

func Test_logEvent_Minimal(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()

	// Call helper
	logEvent(context.Background(), "test-message", AuditCategorySuccess, nil)

	// Assert result
	require.Len(t, hook.Entries, 1)
	expectedData := log.Fields{
		"category": AuditCategorySuccess,
	}
	assert.Equal(t, expectedData, hook.LastEntry().Data)
	assert.Equal(t, "test-message", hook.LastEntry().Message)
	hook.Reset()
}

func Test_logEvent_InvalidMetadata(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	ctx := microMetadata.Set(context.Background(), metadata.GoMicroMetadataKey, "invalid")

	// Call helper
	logEvent(ctx, "test-message", AuditCategorySuccess, nil)

	// Assert result
	assert.True(t, len(hook.Entries) > 2)
	assert.Contains(t, hook.Entries[0].Message, "Failed to decode")
	expectedData := log.Fields{
		"category": AuditCategorySuccess,
	}
	assert.Equal(t, expectedData, hook.LastEntry().Data)
	assert.Equal(t, "test-message", hook.LastEntry().Message)
	hook.Reset()
}
