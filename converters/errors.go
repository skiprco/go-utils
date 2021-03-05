package converters

// ErrorInputIsNotPointer indicates a pointer was expected as input,
// but the provided input is not a pointer.
const ErrorInputIsNotPointer = "input_is_not_a_pointer"

// ErrorPanicDuringSanitizeObject indicates we recovered from a panic
// while running converters.SanitizeObject
const ErrorPanicDuringSanitizeObject = "panic_during_sanitize_object"
