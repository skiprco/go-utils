package auth

// =====================================
// =               ROLES               =
// =====================================

// RoleUser is a user role which means the user is using the Skipr application
const RoleUser = "USER"

// RoleOperatorRead is a user role means the operator has read-only access to all data
const RoleOperatorRead = "OPERATOR_READ"

// RoleOperatorWrite is a user role means the operator has read-write access to all data
const RoleOperatorWrite = "OPERATOR_WRITE"

// RoleOperatorAdmin is a user role means the user has read-write access to all data and can modify the roles of other users
const RoleOperatorAdmin = "OPERATOR_ADMIN"

// =====================================
// =           ROLE MAPPINGS           =
// =====================================

var roleMapUser = map[string]bool{
	RoleUser: true,
}

var roleMapOperatorRead = map[string]bool{
	RoleOperatorRead:  true,
	RoleOperatorWrite: true,
	RoleOperatorAdmin: true,
}

var roleMapOperatorWrite = map[string]bool{
	RoleOperatorWrite: true,
	RoleOperatorAdmin: true,
}

var roleMapOperatorAdmin = map[string]bool{
	RoleOperatorAdmin: true,
}

// =====================================
// =               OTHER               =
// =====================================

// AuthOverridePrefix contains the prefix you have to use to override the authentication.
// This is needed when the action is not invoked by a user (e.g. callback by provider).
const AuthOverridePrefix = "system_override_"
