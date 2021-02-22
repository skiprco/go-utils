package auth

// ErrorNotEnoughPrivileges indicates the user tried to access a
// resource for which it doesn't have enough privileges.
const ErrorNotEnoughPrivileges = "not_enough_privileges"

// AuthOverridePrefix contains the prefix you have to use to override the authentication.
// This is needed when the action is not invoked by a user (e.g. callback by provider).
const AuthOverridePrefix = "system_override_"

// RoleUser is a user role which means the user is using the Skipr application
const RoleUser = "USER"

// RoleOperatorRead is a user role means the operator has read-only access to all data
const RoleOperatorRead = "OPERATOR_READ"

// RoleOperatorWrite is a user role means the operator has read-write access to all data
const RoleOperatorWrite = "OPERATOR_WRITE"

// RoleOperatorAdmin is a user role means the user has read-write access to all data and can modify the roles of other users
const RoleOperatorAdmin = "OPERATOR_ADMIN"
