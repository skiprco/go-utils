package metadata

import (
	"context"

	"github.com/micro/go-micro/v2/metadata"
	"github.com/skiprco/go-utils/v2/collections"
	"github.com/skiprco/go-utils/v2/errors"
)

// GoMicroMetadataKey is the key which is used to store metadata on the go-micro context
const GoMicroMetadataKey = "Metadata"

// GetGoMicroMetadata returns the currently defined metadata from the go-micro context
func GetGoMicroMetadata(ctx context.Context) (Metadata, *errors.GenericError) {
	// Extract encoded metadata from context
	metaBase64, exists := metadata.Get(ctx, GoMicroMetadataKey)
	if !exists {
		return Metadata{}, nil
	}

	// Decode and return metadata
	return FromBase64(metaBase64)
}

// UpdateGoMicroMetadata upserts the metadata stored in the go-micro context.
// The provided metadata will be merged with the currently defined metadata.
// Returns result of the merge.
func UpdateGoMicroMetadata(ctx context.Context, additionalMetadata Metadata) (context.Context, Metadata, *errors.GenericError) {
	// Get current metadata
	meta, genErr := GetGoMicroMetadata(ctx)
	if genErr != nil {
		return ctx, nil, genErr
	}

	// Update and return result
	meta = collections.StringMapMerge(meta, additionalMetadata)
	ctx = metadata.Set(ctx, GoMicroMetadataKey, meta.ToBase64())
	return ctx, meta, nil
}

// SetGoMicroMetadata upserts a single key/value pair in the go-micro context.
// The provided metadata will be merged with the currently defined metadata.
// Returns result of the merge.
func SetGoMicroMetadata(ctx context.Context, key string, value string) (context.Context, Metadata, *errors.GenericError) {
	meta := Metadata{key: value}
	return UpdateGoMicroMetadata(ctx, meta)
}
