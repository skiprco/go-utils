package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkRoles_HasRole(t *testing.T) {
	roles := []string{RoleUser, RoleOperatorRead}
	allowed := map[string]bool{RoleOperatorRead: true}
	result := checkRoles(allowed, roles)
	assert.True(t, result)
}

func Test_checkRoles_DoesntHasRole(t *testing.T) {
	roles := []string{RoleUser, RoleOperatorRead}
	allowed := map[string]bool{RoleOperatorWrite: true}
	result := checkRoles(allowed, roles)
	assert.False(t, result)
}

func Test_checkRoles_RolesIsNil(t *testing.T) {
	allowed := map[string]bool{RoleOperatorWrite: true}
	result := checkRoles(allowed, nil)
	assert.False(t, result)
}
