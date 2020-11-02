package logging

import (
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

func fixtureRequest() (*httptest.ResponseRecorder, *gin.Context) {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest("GET", "/", nil)
	return rec, ctx
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

func Test_AuditMiddleware_OperationSuccess(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	_, ctx := fixtureRequest()

	// Call middleware
	mw := AuditMiddleware("test-operator")
	mw(ctx)

	// Assert operation attempt
	require.Len(t, hook.Entries, 2)
	actualAttempt := testAssertAndDropDynamicEntryFields(t, hook.Entries[0])
	expected := log.Fields{
		"operator": "test-operator",
		"category": auditCategoryAttempt,
		"message":  auditMessageOperationAttempt,
	}
	assert.Equal(t, expected, actualAttempt)

	// Assert operation success
	actualResult := testAssertAndDropDynamicEntryFields(t, hook.Entries[1])
	expected["category"] = auditCategorySuccess
	expected["message"] = auditMessageOperationSuccess
	assert.Equal(t, expected, actualResult)

	// Assert combined
	assert.Equal(t, hook.Entries[0].Data["operation_start_datetime"], hook.Entries[1].Data["operation_start_datetime"])
	hook.Reset()
}

func Test_AuditMiddleware_OperationFail(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	_, ctx := fixtureRequest()

	// Call middleware
	mw := AuditMiddleware("test-operator")
	ctx.Status(400)
	mw(ctx)

	// Assert operation attempt
	require.Len(t, hook.Entries, 2)
	actualAttempt := testAssertAndDropDynamicEntryFields(t, hook.Entries[0])
	expected := log.Fields{
		"operator": "test-operator",
		"category": auditCategoryAttempt,
		"message":  auditMessageOperationAttempt,
	}
	assert.Equal(t, expected, actualAttempt)

	// Assert operation fail
	actualResult := testAssertAndDropDynamicEntryFields(t, hook.Entries[1])
	expected["category"] = auditCategoryFail
	expected["message"] = auditMessageOperationFail
	assert.Equal(t, expected, actualResult)

	// Assert combined
	assert.Equal(t, hook.Entries[0].Data["operation_start_datetime"], hook.Entries[1].Data["operation_start_datetime"])
	hook.Reset()
}
