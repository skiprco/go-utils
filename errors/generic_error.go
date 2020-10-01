package errors

import (
	"fmt"
	"strings"

	microErrors "github.com/micro/go-micro/v2/errors"
)

// GenericError provides a structured output for easy use in the front end
type GenericError struct {
	ID            string
	Code          int
	Status        string
	Domain        string
	SubDomain     string
	SubDomainCode string
	Meta          map[string]string
	IsLegacyError bool
}

// GetDetailString returns the detail string including meta data
func (e GenericError) GetDetailString() string {
	// Check if legacy error
	if e.IsLegacyError {
		return e.Domain
	}

	// Build meta string
	metaList := []string{}
	for key, value := range e.Meta {
		metaList = append(metaList, fmt.Sprintf("%s=%s", key, value))
	}

	// Append meta string
	detailString := e.Domain + "/" + e.SubDomain + "/" + e.SubDomainCode + "/"
	if len(metaList) > 0 {
		metaString := strings.Join(metaList, ";")
		detailString += metaString
	}

	// Return result
	return detailString
}

// Error converts the error into a string
func (e GenericError) Error() string {
	return fmt.Sprintf(
		`{"id": %s, "code": %d, "detail": %s, "status": %s}`,
		e.ID, e.Code, e.GetDetailString(), e.Status,
	)
}

// ToMicroError converts to generic error to a micro error
func (e GenericError) ToMicroError() error {
	return microErrors.New(e.ID, e.GetDetailString(), int32(e.Code))
}
