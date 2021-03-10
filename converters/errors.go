package converters

const errorDomain = "go_utils"
const errorSubDomain = "converters"

// ErrorCountryNotFound indicates the specified country is not found.
const ErrorCountryNotFound = "country_not_found"

// ErrorInputIsNotPointer indicates a pointer was expected as input,
// but the provided input is not a pointer.
const ErrorInputIsNotPointer = "input_is_not_a_pointer"

// ErrorPanicDuringSanitizeObject indicates we recovered from a panic
// while running converters.SanitizeObject.
const ErrorPanicDuringSanitizeObject = "panic_during_sanitize_object"

// ErrorFailedToNormaliseString indicates we failed to normalise the
// provided string. More info is printed in the logs.
const ErrorFailedToNormaliseString = "failed_to_normalise_string"

// ErrorUnknownCurrency indicates the expected currency was not found
const ErrorUnknownCurrency = "unknown_currency"

// ErrorUnknownCurrency indicates the date param cannot be converted.
// Probably related to a wrong format
const ErrorCannotConvertDate = "cannot_convert_date"
