package gin

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"github.com/skiprco/go-utils/v3/logging"
	"github.com/skiprco/go-utils/v3/metadata"
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
		auditData := metadata.Metadata{
			"operator":                 operator,
			"operation_id":             uuid.New(),
			"operation_start_datetime": time.Now().UTC().Format(time.RFC3339),
			"operation_time":           "0",
		}

		// Update metadataToPass
		metadata.UpdateGinMetadata(c, auditData)

		// Log operation attempt
		ctx := metadata.ConvertGinToGoMicro(c)
		additional := map[string]interface{}{ // This data shouldn't be included in the context
			"request_payload": string(bodyBytes),
		}
		logging.AuditOperationAttempt(ctx, additional)

		// Process api call
		c.Next()

		// Read response body
		additional["response_payload"] = blw.body.String()

		// Log operation result
		ctx = metadata.ConvertGinToGoMicro(c)
		if c.Writer.Status() < 300 {
			logging.AuditOperationSuccess(ctx, additional)
		} else {
			logging.AuditOperationFail(ctx, additional)
		}
	}
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
