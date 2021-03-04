package converters

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/skiprco/go-utils/v2/errors"
)

// Sanitize removes all HTML tags from the input and escapes entities.
// Following entities are excluded from escaping: &, ' (apos)
func Sanitize(input string) string {
	// Sanitize input
	output := bluemonday.StrictPolicy().Sanitize(input)

	// Restore characters which shouldn't pose a threat
	output = strings.ReplaceAll(output, "&#39;", "'")
	return strings.ReplaceAll(output, "&amp;", "&")
}

// SanitizeObject takes a pointer to an object (struct, map, slice, ...) as input
// and runs Sanitize for each field which is a (pointer to a) string.
//
// Raises
//
// 500/input_is_not_a_pointer: Provided input is not a pointer
//
// 500/panic_during_sanitize_object: A panic occured during sanitation
func SanitizeObject(input interface{}) (genErr *errors.GenericError) {
	// Validate type of input
	inputType := reflect.TypeOf(input)
	if inputType.Kind() != reflect.Ptr {
		meta := map[string]string{"type": inputType.String()}
		return errors.NewGenericError(500, "go-utils", "common", "input_is_not_a_pointer", meta)
	}

	// Convert panic to correct error
	defer func() {
		if r := recover(); r != nil {
			meta := map[string]string{"panic": fmt.Sprintf("%v", r)}
			genErr = errors.NewGenericError(500, "go-utils", "common", "panic_during_sanitize_object", meta)
		}
	}()

	// Start traverse
	inputValue := reflect.ValueOf(input).Elem()
	sanitizeTraverse(inputValue)

	// Sanitize successful
	return nil
}

func sanitizeTraverse(input reflect.Value) {
	// Unpack if pointer or interface
	if input.Kind() == reflect.Interface || input.Kind() == reflect.Ptr {
		input = input.Elem()
	}

	// Check type and proceed traverse
	switch input.Kind() {
	case reflect.String:
		input.SetString(Sanitize(input.String()))
	}
}
