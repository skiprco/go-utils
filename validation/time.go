package validation

import (
	"time"

	"github.com/skiprco/go-utils/v2/errors"
)

// WithinTimeRange checks if the "nowTime" lies between "startTime" (including) and "endTime" (excluding).
//
// Raises
//
// - 400/end_time_before_start_time: Provided endTime lies before startTime
func WithinTimeRange(nowTime time.Time, startTime time.Time, endTime time.Time) (bool, *errors.GenericError) {
	// endTime should be after startTime
	if endTime.Before(startTime) {
		meta := map[string]string{
			"start_time": startTime.Format(time.RFC3339),
			"end_time":   endTime.Format(time.RFC3339),
		}
		return false, errors.NewGenericError(400, "go-utils", "common", "end_time_before_start_time", meta)
	}

	// Check if within time range
	return nowTime == startTime || (nowTime.After(startTime) && nowTime.Before(endTime)), nil
}
