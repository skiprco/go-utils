package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/skiprco/go-utils/v2/errors"
	"github.com/skiprco/go-utils/v2/logging"
)

// IsOverride checks if the provided user ID has a prefix to override the authentication.
// Overrides will be clearly logged.
func IsOverride(ctx context.Context, userID string, subDomain string) bool {
	if strings.HasPrefix(userID, AuthOverridePrefix) {
		overrideBy := strings.TrimPrefix(userID, AuthOverridePrefix)
		message := fmt.Sprintf("Authorization override for %s by %s", subDomain, overrideBy)
		meta := map[string]interface{}{"override_by": overrideBy}
		logging.AuditFact(ctx, message, meta)
		return true
	}
	return false
}

// HasRoleCheck is a helper function to validate if a user has a specific role
type HasRoleCheck func(ctx context.Context, userID string) (bool, *errors.GenericError)

// MustHaveRole is a helper to ease the implementation of access control checks.
// This helper should be called by a specific helper for the service which implements the access control.
func MustHaveRole(ctx context.Context, hasRoleCheck HasRoleCheck, userID string, errorDomain string, subDomain string) *errors.GenericError {
	// Check for override
	if IsOverride(ctx, userID, subDomain) {
		return nil
	}

	// Not an override => Check roles
	hasRole, genErr := hasRoleCheck(ctx, userID)
	if genErr != nil {
		return genErr
	}
	if !hasRole {
		return errors.NewGenericError(403, errorDomain, subDomain, ErrorNotEnoughPrivileges, nil)
	}
	return nil
}
