package logging

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/pborman/uuid"
	"github.com/skiprco/go-utils/collections"
)

// AuditMiddleware logs the attempt and result for an API call.
// It also sets metadata in the context to support further audit logging.
//
// Arguments
//
// - operator: e.g. booking-api, registration-api, ...
func AuditMiddleware(operator string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build audit data
		auditData := map[string]string{
			"operator":                 operator,
			"operation_id":             uuid.New(),
			"operation_start_datetime": time.Now().UTC().Format(time.RFC3339),
			"operation_time":           "0",
		}

		// Update metadataToPass
		metadata := c.GetStringMapString("metadataToPass")
		metadata = collections.StringMapMerge(metadata, auditData)
		c.Set("metadataToPass", metadata)

		// Log operation attempt
		ctx := getContextFromGin(c)
		AuditOperationAttempt(ctx, nil)

		// Process api call
		c.Next()

		// Log operation result
		ctx = getContextFromGin(c)
		if c.Writer.Status() < 300 {
			AuditOperationSuccess(ctx, nil)
		} else {
			AuditOperationFail(ctx, nil)
		}
	}
}

func getContextFromGin(c *gin.Context) context.Context {
	// Extract metadata map from Gin
	meta := c.GetStringMapString("metadataToPass")

	// Inject each key-value pair in the context
	ctx := c.Request.Context()
	for k, v := range meta {
		ctx = metadata.Set(ctx, k, v)
	}

	// Return result
	return ctx
}
