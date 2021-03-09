package errors

import (
	"testing"

	microErrors "github.com/asim/go-micro/v3/errors"
	"github.com/stretchr/testify/assert"
)

func testCopyMicroToGenericErrorFields(microError *microErrors.Error) *GenericError {
	return &GenericError{
		ID:            microError.Id,
		Code:          int(microError.Code),
		Status:        microError.Status,
		Meta:          map[string]string{},
		IsLegacyError: false,
	}
}

func TestNewGenericFromMicroError(t *testing.T) {
	// Invalid error
	invalidMicro := microErrors.Parse("invalid-error")
	invalidGeneric := testCopyMicroToGenericErrorFields(invalidMicro)
	invalidGeneric.Domain = "invalid-error"
	invalidGeneric.IsLegacyError = true

	// Error without meta
	withoutMetaMicro := microErrors.Parse("domain/subdomain/error/")
	withoutMetaGeneric := testCopyMicroToGenericErrorFields(withoutMetaMicro)
	withoutMetaGeneric.Domain = "domain"
	withoutMetaGeneric.SubDomain = "subdomain"
	withoutMetaGeneric.SubDomainCode = "error"

	// Error with meta
	withMetaMicro := microErrors.Parse("domain/subdomain/error/test=success;test2=failed")
	withMetaGeneric := testCopyMicroToGenericErrorFields(withMetaMicro)
	withMetaGeneric.Domain = "domain"
	withMetaGeneric.SubDomain = "subdomain"
	withMetaGeneric.SubDomainCode = "error"
	withMetaGeneric.Meta = map[string]string{
		"test":  "success",
		"test2": "failed",
	}

	// Error with invalid meta
	invalidMetaMicro := microErrors.Parse("domain/subdomain/error/test=success;test2-failed;test3=fixed")
	invalidMetaGeneric := testCopyMicroToGenericErrorFields(invalidMetaMicro)
	invalidMetaGeneric.Domain = "domain"
	invalidMetaGeneric.SubDomain = "subdomain"
	invalidMetaGeneric.SubDomainCode = "error"
	invalidMetaGeneric.Meta = map[string]string{
		"test":  "success",
		"test3": "fixed",
	}

	// Error with invalid meta converted to legacy
	invalidMetaMicroLegacy := microErrors.Parse("domain/subdomain/error/test=success;test2/failed;test3=fixed")
	invalidMetaGenericLegacy := testCopyMicroToGenericErrorFields(invalidMetaMicro)
	invalidMetaGenericLegacy.Domain = "domain/subdomain/error/test=success;test2/failed;test3=fixed"
	invalidMetaGenericLegacy.IsLegacyError = true

	expectations := map[error]*GenericError{
		invalidMicro:           invalidGeneric,
		withoutMetaMicro:       withoutMetaGeneric,
		withMetaMicro:          withMetaGeneric,
		invalidMetaMicro:       invalidMetaGeneric,
		invalidMetaMicroLegacy: invalidMetaGenericLegacy,
	}

	for input, expected := range expectations {
		t.Run("convert micro to generic error", func(t *testing.T) {
			result := NewGenericFromMicroError(input)
			assert.Equal(t, expected, result)
		})
	}
}
