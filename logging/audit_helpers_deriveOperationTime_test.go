package logging

import (
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_deriveOperationTime_Empty_Success(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	testFields := logrus.Fields{}

	// Call helper
	deriveOperationTime(testFields)

	// Assert result
	assert.Len(t, hook.Entries, 0)
	assert.Equal(t, logrus.Fields{}, testFields)
	hook.Reset()
}

func Test_deriveOperationTime_Specified_Success(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	testFields := logrus.Fields{
		"operation_start_datetime": time.Now().Add(-time.Hour).Format(time.RFC3339),
	}

	// Call helper
	deriveOperationTime(testFields)

	// Assert result
	assert.Len(t, hook.Entries, 0)
	operationTime, err := strconv.Atoi(testFields["operation_time"].(string))
	require.Nil(t, err)
	assert.InDelta(t, operationTime, time.Hour.Milliseconds(), float64(time.Second.Milliseconds())) // Check within 1 second accuracy
	hook.Reset()
}

func Test_deriveOperationTime_InvalidType_Failure(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	testFields := logrus.Fields{
		"operation_start_datetime": time.Now(),
	}

	// Call helper
	deriveOperationTime(testFields)

	// Assert result
	assert.Len(t, hook.Entries, 1)
	hook.Reset()
}

func Test_deriveOperationTime_InvalidFormat_Failure(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	testFields := logrus.Fields{
		"operation_start_datetime": time.Now().Format(time.Kitchen),
	}

	// Call helper
	deriveOperationTime(testFields)

	// Assert result
	assert.Len(t, hook.Entries, 1)
	hook.Reset()
}

func Test_deriveOperationTime_ZeroDate_Failure(t *testing.T) {
	// Setup test
	hook := logTest.NewGlobal()
	testFields := logrus.Fields{
		"operation_start_datetime": time.Time{}.Format(time.RFC3339),
	}

	// Call helper
	deriveOperationTime(testFields)

	// Assert result
	assert.Len(t, hook.Entries, 1)
	hook.Reset()
}
