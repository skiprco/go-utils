package logging

import (
	"context"
	"errors"
	"testing"

	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/server"
	"github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MicroRequest struct{}

func (req MicroRequest) Service() string           { return "test-request-service" }
func (req MicroRequest) Method() string            { return "test-request-method" }
func (req MicroRequest) Endpoint() string          { return "test-request-endpoint" }
func (req MicroRequest) ContentType() string       { return "test-request-content-type" }
func (req MicroRequest) Header() map[string]string { return nil }
func (req MicroRequest) Body() interface{}         { return nil }
func (req MicroRequest) Read() ([]byte, error)     { return nil, nil }
func (req MicroRequest) Codec() codec.Reader       { return nil }
func (req MicroRequest) Stream() bool              { return false }

func testAssertAuditHandlerWrapperEntry(t *testing.T, logEntry logrus.Entry, category AuditCategory) {
	// Assert log fields
	assert.Equal(t, "test-request-service", logEntry.Data["service_name"])
	assert.Equal(t, "test-request-endpoint", logEntry.Data["service_endpoint"])
	assert.EqualValues(t, category, logEntry.Data["category"])

	// Assert message
	assert.Equal(t, "test-request-service_test-request-endpoint", logEntry.Message)
}

func Test_AuditHandlerWrapper_Success(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	handler := func(ctx context.Context, req server.Request, rsp interface{}) error { return nil }
	request := MicroRequest{}

	// Call helper
	wrapper := AuditHandlerWrapper(handler)
	err := wrapper(context.Background(), request, nil)

	// Assert result
	assert.Nil(t, err)
	require.Len(t, hook.Entries, 2)
	testAssertAuditHandlerWrapperEntry(t, hook.Entries[0], AuditCategoryAttempt)
	testAssertAuditHandlerWrapperEntry(t, hook.Entries[1], AuditCategorySuccess)
	hook.Reset()
}

func Test_AuditHandlerWrapper_Failuer(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	handler := func(ctx context.Context, req server.Request, rsp interface{}) error { return errors.New("test-error") }
	request := MicroRequest{}

	// Call helper
	wrapper := AuditHandlerWrapper(handler)
	err := wrapper(context.Background(), request, nil)

	// Assert result
	assert.EqualError(t, err, "test-error")
	require.Len(t, hook.Entries, 2)
	testAssertAuditHandlerWrapperEntry(t, hook.Entries[0], AuditCategoryAttempt)
	testAssertAuditHandlerWrapperEntry(t, hook.Entries[1], AuditCategoryFail)
	hook.Reset()
}
