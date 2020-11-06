package metadata

import (
	"context"

	"github.com/gin-gonic/gin"
)

// ConvertGinToGoMicro returns a context object with the metadata
// set as go-micro metadata. This way the metadata can be accessed in
// each microservices.
func ConvertGinToGoMicro(c *gin.Context) context.Context {
	// Get base context from request
	ctx := c.Request.Context()

	// "genErr" ignored since it can only fail on extracting metadata from context.
	// This cannot occur since there is no go-micro metadata on the context yet.
	meta := GetGinMetadata(c)
	ctx, _, _ = UpdateGoMicroMetadata(ctx, meta)

	// Return result
	return ctx
}
