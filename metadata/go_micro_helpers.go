package metadata

import (
	"context"

	"github.com/skiprco/go-utils/v2/errors"
)

// GetUserIDFromGoMicroMeta extracts the user ID from the metadata of go-micro
// and also returns the raw metadata for later use. Throws an error if unable
// to read metadata or if user_id is not set.
//
// Raises
//
// - 500/user_id_not_set_in_metadata: The user ID key is not set in the metadata
func GetUserIDFromGoMicroMeta(ctx context.Context, errorDomain string) (string, Metadata, *errors.GenericError) {
	// Get user from metadata
	meta, genErr := GetGoMicroMetadata(ctx)
	if genErr != nil {
		return "", meta, genErr
	}

	// Validate if user ID is set
	userID := meta.Get("user_id")
	if userID == "" {
		return "", meta, errors.NewGenericError(500, errorDomain, errorSubDomain, ErrorUserIDNotInMeta, nil)
	}

	return userID, meta, nil
}
