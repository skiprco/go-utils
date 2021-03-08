package converters

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/v2/errors"
)

// Sanitize removes all HTML tags from the input and escapes entities.
// Following entities are excluded from escaping: ' (apos), " (quote), &
func Sanitize(input string) string {
	// Sanitize input
	output := bluemonday.StrictPolicy().Sanitize(input)

	// Restore some characters
	output = strings.ReplaceAll(output, "&#39;", "'")
	output = strings.ReplaceAll(output, "&#34;", `"`)
	return strings.ReplaceAll(output, "&amp;", "&")
}

// SanitizeObject takes a pointer to an object (struct, map, slice, ...) as input
// and runs Sanitize for each field which is a (pointer to a) string.
//
// Raises
//
// - 500/input_is_not_a_pointer: Provided input is not a pointer
//
// - 500/panic_during_sanitize_object: A panic occured during sanitation
func SanitizeObject(input interface{}) (genErr *errors.GenericError) {
	// Validate type of input
	inputType := reflect.TypeOf(input)
	if inputType.Kind() != reflect.Ptr {
		log.WithField("input", input).Error("Provided input to sanitize must be a pointer")
		meta := map[string]string{"type": inputType.String()}
		return errors.NewGenericError(500, "go-utils", "common", ErrorInputIsNotPointer, meta)
	}

	// Convert panic to correct error
	defer func() {
		if r := recover(); r != nil {
			meta := map[string]string{"panic": fmt.Sprintf("%v", r)}
			genErr = errors.NewGenericError(500, "go-utils", "common", ErrorPanicDuringSanitizeObject, meta)
			log.WithFields(log.Fields{
				"error": genErr,
				"input": input,
			}).Error("Panic thrown during sanitation")
		}
	}()

	// Start traverse
	inputValue := reflect.ValueOf(input).Elem()
	sanitizeTraverse(inputValue)

	// Sanitize successful
	return nil
}

// sanitizeTraverse is a recursive function which sanitizes each child of the provided input
func sanitizeTraverse(input reflect.Value) *errors.GenericError {
	// Unpack if interface or pointer
	if input.Kind() == reflect.Interface || input.Kind() == reflect.Ptr {
		input = input.Elem()
	}

	// Check type and proceed traverse
	switch input.Kind() {
	case reflect.String:
		input.SetString(Sanitize(input.String()))

	case reflect.Struct:
		for i := 0; i < input.NumField(); i++ {
			name := input.Type().Field(i).Name
			if isExported(name) {
				// Field is exported => Continue traversal
				genErr := sanitizeTraverse(input.Field(i))
				if genErr != nil {
					return genErr
				}
			}
		}

	case reflect.Map:
		genErr := sanitizeMap(input)
		if genErr != nil {
			return genErr
		}

	case reflect.Slice:
		for i := 0; i < input.Len(); i++ {
			genErr := sanitizeTraverse(input.Index(i))
			if genErr != nil {
				return genErr
			}
		}
	}

	// Traverse succesful
	return nil
}

func sanitizeMap(input reflect.Value) *errors.GenericError {
	for _, k := range input.MapKeys() {
		// Clean copy of key and value
		cleanedKey := sanitizeCopy(k)
		value := input.MapIndex(k)
		cleanedValue := sanitizeCopy(value)

		// Delete old key from map
		input.SetMapIndex(k, reflect.Value{})

		// Write new key/value to map
		input.SetMapIndex(cleanedKey, cleanedValue)
	}

	// Sanitize successful
	return nil
}

// sanitizeCopy takes a Value which has "CanSet() == false",
// makes a copy, sanitizes the copy and returns this copy.
//
// This is needed when you want to sanitize e.g. a value from a map.
// See https://golang.org/pkg/reflect/#Value.CanSet for more info.
func sanitizeCopy(input reflect.Value) reflect.Value {
	var cleaned reflect.Value
	if input.Kind() == reflect.Ptr {
		cleaned = reflect.New(input.Type())
		if !input.IsNil() {
			cleanedValue := reflect.New(input.Elem().Type())
			cleanedValue.Elem().Set(input.Elem())
			sanitizeTraverse(cleanedValue)
			cleaned.Elem().Set(cleanedValue)
		}
	} else {
		cleaned = reflect.New(input.Type())
		cleaned.Elem().Set(input)
		sanitizeTraverse(cleaned)
	}
	return cleaned.Elem()
}

// isExported checks if the provided field name is exported
// by checking if the first letter is capital
func isExported(input string) bool {
	// Based on https://stackoverflow.com/a/30263910
	for _, rune := range input {
		return unicode.IsUpper(rune)
	}
	return false
}
