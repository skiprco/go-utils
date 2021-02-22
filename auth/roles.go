package auth

// OperatorRead checks if the provided roles contain the OPERATOR_READ role.
// This role might be granted implicitely (e.g. on OPERATOR_ADMIN).
func OperatorRead(roles []string) bool {
	readRoles := map[string]bool{
		RoleOperatorRead:  true,
		RoleOperatorWrite: true,
		RoleOperatorAdmin: true,
	}
	return checkRoles(readRoles, roles)
}

// OperatorWrite checks if the provided roles contain the OPERATOR_WRITE role.
// This role might be granted implicitely (e.g. on OPERATOR_ADMIN).
func OperatorWrite(roles []string) bool {
	readRoles := map[string]bool{
		RoleOperatorWrite: true,
		RoleOperatorAdmin: true,
	}
	return checkRoles(readRoles, roles)
}

func checkRoles(allowedRoles map[string]bool, userRoles []string) bool {
	for _, role := range userRoles {
		if allowedRoles[role] {
			return true
		}
	}

	return false
}
