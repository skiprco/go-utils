package http

const errorDomain = "go_utils"
const errorSubDomain = "http"

// ErrorMarshalRequestBodyFailed indicates the marshaling of the request body to JSON failed.
const ErrorMarshalRequestBodyFailed = "marshal_request_body_failed"

// ErrorSendHTTPRequestFailed indicates sending the request failed (e.g. server unreachable).
const ErrorSendHTTPRequestFailed = "send_http_request_failed"

// ErrorReadResponseBodyFailed indicates reading the body to
// bytes (response code < 300) or a string (response code >= 300) failed.
const ErrorReadResponseBodyFailed = "read_response_body_failed"

// ErrorResponseCodeIsError indicates the server returned an error response code.
const ErrorResponseCodeIsError = "response_code_is_error"

// ErrorParseResponseBodyFailed indicates parsing the JSON response into provided response interface failed.
const ErrorParseResponseBodyFailed = "parse_response_body_failed"
