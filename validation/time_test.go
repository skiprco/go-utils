package validation

import (
	"testing"
	"time"

	"github.com/skiprco/go-utils/v3/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WithinTimeRange_WithinTimeRange_True_Success(t *testing.T) {
	now := time.Now()
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	within, genErr := WithinTimeRange(now, start, end)
	require.Nil(t, genErr)
	assert.True(t, within)
}

func Test_WithinTimeRange_SameAsStart_True_Success(t *testing.T) {
	now := time.Now()
	start := now
	end := now.Add(time.Hour)
	within, genErr := WithinTimeRange(now, start, end)
	require.Nil(t, genErr)
	assert.True(t, within)
}

func Test_WithinTimeRange_LiesBefore_False_Success(t *testing.T) {
	now := time.Now().Add(-time.Hour)
	start := time.Now()
	end := time.Now().Add(time.Hour)
	within, genErr := WithinTimeRange(now, start, end)
	require.Nil(t, genErr)
	assert.False(t, within)
}

func Test_WithinTimeRange_LiesAfter_False_Success(t *testing.T) {
	now := time.Now().Add(2 * time.Hour)
	start := time.Now()
	end := time.Now().Add(time.Hour)
	within, genErr := WithinTimeRange(now, start, end)
	require.Nil(t, genErr)
	assert.False(t, within)
}

func Test_WithinTimeRange_EndBeforeStart_Failure(t *testing.T) {
	now := time.Now()
	start := now.Add(time.Hour)
	end := now
	within, genErr := WithinTimeRange(now, start, end)
	errors.AssertGenericError(t, genErr, 400, ErrorEndTimeBeforeStartTime, nil)
	assert.False(t, within)
}
