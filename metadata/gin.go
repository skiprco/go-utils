package metadata

import (
	"github.com/gin-gonic/gin"
	"github.com/skiprco/go-utils/v3/collections"
)

// GinMetadataKey is the key which is used to store metadata on the Gin context
const GinMetadataKey string = "metadataToPass"

// GetGinMetadata returns the currently defined metadata from the Gin context
func GetGinMetadata(c *gin.Context) Metadata {
	meta := c.GetStringMapString(GinMetadataKey)
	if meta == nil {
		return Metadata{}
	}
	return meta
}

// UpdateGinMetadata upserts the metadata stored in the Gin context.
// The provided metadata will be merged with the currently defined metadata.
// Returns result of the merge.
func UpdateGinMetadata(c *gin.Context, additionalMetadata Metadata) Metadata {
	// Merge and update metadata
	currentMetadata := GetGinMetadata(c)
	newMetadata := collections.StringMapMerge(currentMetadata, additionalMetadata)
	c.Set(GinMetadataKey, newMetadata)
	return newMetadata
}

// SetGinMetadata upserts a single key/value pair in the Gin context.
// The provided metadata will be merged with the currently defined metadata.
// Returns result of the merge.
func SetGinMetadata(c *gin.Context, key string, value string) Metadata {
	meta := Metadata{key: value}
	return UpdateGinMetadata(c, meta)
}
