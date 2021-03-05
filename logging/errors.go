package logging

const errorDomain = "go_utils"
const errorSubDomain = "logging"

// ErrorKeyNotFoundInContext indicates we expected to find a key in
// the metadata of the context, but the key is not set.
const ErrorKeyNotFoundInContext = "key_not_found_in_context"

// ErrorUnableToDumpRequest indicates dumping the HTTP request failed
const ErrorUnableToDumpRequest = "unable_to_dump_request"

// ErrorUnableToDumpResponse indicates dumping the HTTP response failed
const ErrorUnableToDumpResponse = "unable_to_dump_response"
