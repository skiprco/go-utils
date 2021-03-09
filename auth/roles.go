package auth

import "github.com/skiprco/go-utils/v3/errors"

// HasRole checks if the provided roles is included in the user roles.
// This role might be granted implicitely (e.g. OPERATOR_READ on OPERATOR_ADMIN).
//
// Raises
//
// - 400/unknown_role: Provided role does not exist
func HasRole(role string, userRoles []string) (bool, *errors.GenericError) {
	switch role {
	case RoleUser:
		return checkRoles(roleMapUser, userRoles), nil
	case RoleOperatorRead:
		return checkRoles(roleMapOperatorRead, userRoles), nil
	case RoleOperatorWrite:
		return checkRoles(roleMapOperatorWrite, userRoles), nil
	case RoleOperatorAdmin:
		return checkRoles(roleMapOperatorAdmin, userRoles), nil
	default:
		meta := map[string]string{"role": role}
		return false, errors.NewGenericError(400, errorDomain, errorSubDomain, ErrorUnknownRole, meta)
	}
}

func checkRoles(allowedRoles map[string]bool, userRoles []string) bool {
	for _, role := range userRoles {
		if allowedRoles[role] {
			return true
		}
	}

	return false
}
