package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/skiprco/go-utils/collections"
)

const contextMetadataKey string = "metadataToPass"

// GetMetadata returns the currently defined metadata from the context
func GetMetadata(c *gin.Context) map[string]string {
	meta := c.GetStringMapString(contextMetadataKey)
	if meta == nil {
		return map[string]string{}
	}
	return meta
}

// UpdateMetadata upserts the metadata stored in the context.
// The provided metadata will be merged with the currently defined metadata.
// Returns result of the merge.
func UpdateMetadata(c *gin.Context, additionalMetadata map[string]string) map[string]string {
	currentMetadata := GetMetadata(c)
	newMetadata := collections.StringMapMerge(currentMetadata, additionalMetadata)
	c.Set(contextMetadataKey, newMetadata)
	return newMetadata
}

// GetContextWithMetadata returns a context object with the metadata
// set as go-micro metadata. This way the metadata can be accessed in
// each microservices.
func GetContextWithMetadata(c *gin.Context) context.Context {
	// Get base context from request
	ctx := c.Request.Context()

	// Inject each key-value pair in the context
	for k, v := range GetMetadata(c) {
		ctx = metadata.Set(ctx, k, v)
	}

	// Return result
	return ctx
}
