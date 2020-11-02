package logging

import (
	"context"
	"testing"

	"github.com/micro/go-micro/v2/metadata"
	log "github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_logEvent_WithContextAndAdditionalData_Success(t *testing.T) {
	// This test checks if the log data is correctly
	// merged from the different sources (context, parameters, ...)

	// Setup test
	hook := logTest.NewGlobal()
	meta := metadata.Metadata{
		"test-meta-key":          "test-meta-value",
		"test-additional-string": "should-be-overwritten",
	}
	ctx := metadata.NewContext(context.Background(), meta)
	additional := map[string]interface{}{
		"test-additional-string": "test-additional-value",
		"test-additional-int":    123,
		"category":               "should-be-overwritten",
	}

	// Call helper
	logEvent(ctx, "test-message", AuditCategorySuccess, additional)

	// Assert result
	require.Len(t, hook.Entries, 1)
	expectedData := log.Fields{
		"category":               AuditCategorySuccess,
		"test-additional-string": "test-additional-value",
		"test-additional-int":    123,
		"test-meta-key":          "test-meta-value",
	}
	assert.Equal(t, expectedData, hook.LastEntry().Data)
	assert.Equal(t, "test-message", hook.LastEntry().Message)
	hook.Reset()
}

func Test_logEvent_Minimal_Success(t *testing.T) {
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
