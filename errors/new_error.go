package errors

import (
	"net/http"
	"strings"
	"time"

	microErrors "github.com/micro/go-micro/v2/errors"
	log "github.com/sirupsen/logrus"
)

// NewGenericError creates a new generic error.
// WARNING: This function returns a GenericError. Do not assign it to an interface of type error!
func NewGenericError(code int, domain string, subDomain string, subDomainCode string, additionalMeta map[string]string) *GenericError {
	// Build and return error
	return &GenericError{
		ID:            time.Now().UTC().Format(time.RFC3339),
		Code:          code,
		Status:        http.StatusText(code),
		Domain:        domain,
		SubDomain:     subDomain,
		SubDomainCode: subDomainCode,
		Meta:          mergeMeta(defaultMeta, additionalMeta),
		IsLegacyError: false,
	}
}

// NewGenericFromMicroError converts a micro error to a generic error
// WARNING: This function returns a GenericError. Do not assign it to an interface of type error!
func NewGenericFromMicroError(err error) *GenericError {
	// Parse settings
	detailSeparator := "/"
	metaSeparator := ";"

	// Parse error into micro error
	microErr := microErrors.Parse(err.Error())

	// Create base error
	genErr := &GenericError{
		ID:            microErr.Id,
		Code:          int(microErr.Code),
		Status:        microErr.Status,
		Meta:          map[string]string{},
		IsLegacyError: false,
	}

	// Parse detail
	detailParts := strings.Split(microErr.Detail, detailSeparator)
	if len(detailParts) != 4 {
		// Legacy error
		genErr.IsLegacyError = true
		genErr.Domain = microErr.Detail
		log.WithField("error", microErr).Warn("MicroError received with legacy format")
		return genErr
	}

	// Split detail
	genErr.Domain = detailParts[0]
	genErr.SubDomain = detailParts[1]
	genErr.SubDomainCode = detailParts[2]

	// Parse meta
	metaString := detailParts[3]
	if metaString != "" {
		// Extract meta pairs from details
		metaPairs := strings.Split(metaString, metaSeparator)
		for _, metaPair := range metaPairs {
			// Split meta item
			metaItem := strings.Split(metaPair, "=")

			if len(metaItem) != 2 {
				// Skip invalid meta
				continue
			}

			// Append meta
			metaKey := metaItem[0]
			metaValue := metaItem[1]
			genErr.Meta[metaKey] = metaValue
		}
	}

	// Return result
	return genErr
}
