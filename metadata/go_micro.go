package metadata

import (
	"context"

	"github.com/asim/go-micro/v3/metadata"
	"github.com/skiprco/go-utils/v3/collections"
	"github.com/skiprco/go-utils/v3/errors"
)

// GoMicroMetadataKey is the key which is used to store metadata on the go-micro context
const GoMicroMetadataKey = "Metadata"

// GetGoMicroMetadata returns the currently defined metadata from the go-micro context
//
// Raises
//
// - 400/decode_glob_from_base64_failed: Failed to decode the metadata as glob from the provided base64 string
// (Note: This error can only occur if the metadata in the go-micro context got corrupted. Therefore this
// error can never occur on a newly created context)
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
//
// Raises
//
// - 400/decode_glob_from_base64_failed: Failed to decode the metadata as glob from the provided base64 string
// (Note: This error can only occur if the metadata in the go-micro context got corrupted. Therefore this
// error can never occur on a newly created context)
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
//
// Raises
//
// - 400/decode_glob_from_base64_failed: Failed to decode the metadata as glob from the provided base64 string
// (Note: This error can only occur if the metadata in the go-micro context got corrupted. Therefore this
// error can never occur on a newly created context)
func SetGoMicroMetadata(ctx context.Context, key string, value string) (context.Context, Metadata, *errors.GenericError) {
	meta := Metadata{key: value}
	return UpdateGoMicroMetadata(ctx, meta)
}
