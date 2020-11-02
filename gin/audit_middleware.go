package gin

import (
	"bytes"
	"context"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/pborman/uuid"
	"github.com/skiprco/go-utils/collections"
	"github.com/skiprco/go-utils/logging"
)

// AuditMiddleware logs the attempt and result for an API call.
// It also sets metadata in the context to support further audit logging.
//
// Arguments
//
// - operator: e.g. booking-api, registration-api, ...
func AuditMiddleware(operator string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Replace writer with bodyLogWriter to capture response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Build audit data
		auditData := map[string]string{
			"operator":                 operator,
			"operation_id":             uuid.New(),
			"operation_start_datetime": time.Now().UTC().Format(time.RFC3339),
			"operation_time":           "0",
		}

		// Update metadataToPass
		meta := c.GetStringMapString("metadataToPass")
		meta = collections.StringMapMerge(meta, auditData)
		c.Set("metadataToPass", meta)

		// Log operation attempt
		ctx := getContextFromGin(c)
		additional := map[string]interface{}{ // This data shouldn't be included in the context
			"request_payload": string(bodyBytes),
		}
		logging.AuditOperationAttempt(ctx, additional)

		// Process api call
		c.Next()

		// Read response body
		additional["response_payload"] = blw.body.String()

		// Log operation result
		ctx = getContextFromGin(c)
		if c.Writer.Status() < 300 {
			logging.AuditOperationSuccess(ctx, additional)
		} else {
			logging.AuditOperationFail(ctx, additional)
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

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.Write([]byte(s))
	return w.ResponseWriter.WriteString(s)
}
