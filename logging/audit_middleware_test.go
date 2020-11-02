package logging

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// =                 TESTS                =
// ========================================

func Test_AuditMiddleware_OperationSuccess(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	router := fixtureRouter(200, 10*time.Millisecond)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, fixtureRequest())

	// Assert operation attempt
	require.Len(t, hook.Entries, 2)
	actualAttempt := testAssertAndDropDynamicEntryFields(t, hook.Entries[0])
	expected := log.Fields{
		"operator":        "test-operator",
		"category":        auditCategoryAttempt,
		"message":         auditMessageOperationAttempt,
		"request_payload": fixtureBody,
	}
	assert.Equal(t, expected, actualAttempt)

	// Assert operation success
	actualResult := testAssertAndDropDynamicEntryFields(t, hook.Entries[1])
	expected["category"] = auditCategorySuccess
	expected["message"] = auditMessageOperationSuccess
	expected["response_payload"] = "test-response-body"
	assert.Equal(t, expected, actualResult)

	// Assert combined
	assert.Equal(t, hook.Entries[0].Data["operation_start_datetime"], hook.Entries[1].Data["operation_start_datetime"])
	assert.NotEqual(t, hook.Entries[0].Data["operation_time"], hook.Entries[1].Data["operation_time"])
	hook.Reset()
}

func Test_AuditMiddleware_OperationFail(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	router := fixtureRouter(400, 0)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, fixtureRequest())

	// Assert operation attempt
	require.Len(t, hook.Entries, 2)
	actualAttempt := testAssertAndDropDynamicEntryFields(t, hook.Entries[0])
	expected := log.Fields{
		"operator":        "test-operator",
		"category":        auditCategoryAttempt,
		"message":         auditMessageOperationAttempt,
		"request_payload": fixtureBody,
	}
	assert.Equal(t, expected, actualAttempt)

	// Assert operation fail
	actualResult := testAssertAndDropDynamicEntryFields(t, hook.Entries[1])
	expected["category"] = auditCategoryFail
	expected["message"] = auditMessageOperationFail
	expected["response_payload"] = "test-response-body"
	assert.Equal(t, expected, actualResult)

	// Assert combined
	assert.Equal(t, hook.Entries[0].Data["operation_start_datetime"], hook.Entries[1].Data["operation_start_datetime"])
	hook.Reset()
}

// ========================================
// =                HELPERS               =
// ========================================

const fixtureBody = `{"test-key": "test-value"}`

func fixtureRequest() *http.Request {
	bodyReader := ioutil.NopCloser(bytes.NewBufferString(fixtureBody))
	return httptest.NewRequest("POST", "/", bodyReader)
}

func fixtureRouter(status int, duration time.Duration) *gin.Engine {
	// Create router
	router := gin.Default()

	// Set audit middleware
	auditMW := AuditMiddleware("test-operator")
	router.Use(auditMW)

	// Add simple handler
	handler := func(c *gin.Context) {
		time.Sleep(duration) // Some processing
		c.String(status, "test-response-body")
	}
	router.POST("/", handler)

	// Return engine
	return router
}

func testAssertAndDropDynamicEntryFields(t *testing.T, entry log.Entry) log.Fields {
	// Duplicate fields
	fields := make(log.Fields, len(entry.Data))
	for k, v := range entry.Data {
		fields[k] = v
	}

	// Assert operation_id => Valid UUID
	require.Contains(t, fields, "operation_id")
	require.NotNil(t, uuid.Parse(fields["operation_id"].(string)), "operation_id should be a valid UUID")
	delete(fields, "operation_id")

	// Assert operation_start_datetime => Valid RFC3339 and max 2 seconds old
	require.Contains(t, fields, "operation_start_datetime")
	start, err := time.Parse(time.RFC3339, fields["operation_start_datetime"].(string))
	require.Nil(t, err)
	require.WithinDuration(t, time.Now(), start, 2*time.Second)
	delete(fields, "operation_start_datetime")

	// Assert operation_time
	require.Contains(t, fields, "operation_time")
	operationTime, err := strconv.Atoi(fields["operation_time"].(string))
	require.Nil(t, err)
	require.True(t, int64(operationTime) < 2*time.Second.Milliseconds(), "operation_time should be less than 2 seconds")
	delete(fields, "operation_time")

	// Return remaining fields
	return fields
}
