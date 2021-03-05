package auth

const errorDomain = "go_utils"
const errorSubDomain = "auth"

// ErrorNotEnoughPrivileges indicates the user tried to access a
// resource for which it doesn't have enough privileges.
const ErrorNotEnoughPrivileges = "not_enough_privileges"

// ErrorUnknownRole indicates the checked role doesn't exist.
const ErrorUnknownRole = "unknown_role"
