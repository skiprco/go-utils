package logging

import (
	"github.com/gin-gonic/gin"
)

// LogAudit logs the attempt and result for an API call.
// It also sets metadata in the context to support further audit logging.
//
// Arguments
//
// - operator: e.g. booking-api, registration-api, ...
func LogAudit(operator string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build audit data
		auditData := map[string]string{
			"operator": operator,
		}

		// Update metadataToPass
		metadata := c.GetStringMapString("metadataToPass")
	}
}
