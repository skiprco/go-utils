package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/v3/errors"
	"github.com/skiprco/go-utils/v3/logging"
)

// Call marshals the body to JSON and sends a HTTP request. Response is parsed as JSON.
//
// Raises
//
// - 500/marshal_request_body_failed: Failed to marshal request body to JSON
//
// - 421/send_http_request_failed: Failed to send request (e.g. server unreachable)
//
// - 421/read_response_body_failed: Failed to read response body (only tried if response code is < 300)
//
// - 500/read_response_body_failed: Failed to read response body into a string (only tried if response code is >= 300)
//
// - dyn/response_code_is_error: Server returned an error code.
// In case the full response will be returned. Also the response body is present on
// the error meta as key "response_body". This way, you can translate the body to a more
// specific error in an easy way. When you miss a translation, the full body is present on the
// error for easier debugging.
//
// - 421/parse_response_body_failed: Failed to parse JSON response into provided response interface
func Call(method string, rootURL string, path string, body interface{}, response interface{}, query map[string]string, headers map[string]string) (*http.Response, *errors.GenericError) {
	// Set Content-Type to JSON
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/json"

	// Setup default logging fields
	defaultLog := log.WithFields(log.Fields{
		"method":  method,
		"url":     rootURL + path,
		"body":    body,
		"query":   query,
		"headers": headers,
	})

	// Marshal body
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			defaultLog.WithField("error", err).Error("Failed to marshal request body to JSON")
			return nil, errors.NewGenericError(500, errorDomain, errorSubDomain, ErrorMarshalRequestBodyFailed, nil)
		}
	}

	// Create request
	req := createRequest(method, rootURL, path, bodyBytes, query, headers)

	// Send request to API
	res, genErr := sendRequest(req)
	if genErr != nil {
		defaultLog.WithField("error", genErr.GetDetailString()).Error("Failed to send request")
		return res, genErr
	}

	// Debug logging
	traceID, _ := logging.LogHTTPRequestResponse(req, res, log.WarnLevel, "After response")
	traceLog := defaultLog.WithField("trace_id", traceID)

	// Read response body
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		traceLog.WithField("error", genErr.GetDetailString()).Error("Failed to read response body")
		return nil, errors.NewGenericError(421, errorDomain, errorSubDomain, ErrorReadResponseBodyFailed, nil)
	}

	// Unmarshal response
	err = json.Unmarshal(resBody, response)
	if err != nil {
		traceLog.WithField("error", err).Error("Failed to parse response body from JSON")
		return nil, errors.NewGenericError(421, errorDomain, errorSubDomain, ErrorParseResponseBodyFailed, nil)
	}

	// Return result
	return res, nil
}

// CallRaw sends a HTTP request without modifying the body
//
// Raises
//
// - 421/send_http_request_failed: Failed to send request (e.g. server unreachable)
//
// - 500/read_response_body_failed: Failed to read response body into a string (only tried if response code is >= 300)
//
// - dyn/response_code_is_error: Server returned an error code.
// In case the full response will be returned. Also the response body is present on
// the error meta as key "response_body". This way, you can translate the body to a more
// specific error in an easy way. When you miss a translation, the full body is present on the
// error for easier debugging.
func CallRaw(method string, rootURL string, path string, body []byte, query map[string]string, headers map[string]string) (*http.Response, *errors.GenericError) {
	// Create request
	req := createRequest(method, rootURL, path, body, query, headers)

	// Send request to API
	res, genErr := sendRequest(req)
	if genErr != nil {
		log.WithFields(log.Fields{
			"method":  method,
			"url":     rootURL + path,
			"body":    string(body),
			"query":   query,
			"headers": headers,
			"error":   genErr.GetDetailString(),
		}).Error("Failed to send HTTP request")
		return res, genErr
	}
	// Return result
	return res, nil
}

func createRequest(method string, rootURL string, path string, body []byte, query map[string]string, headers map[string]string) *http.Request {
	// Create HTTP request
	req, _ := http.NewRequest(method, rootURL+path, bytes.NewBuffer(body))

	// Apply custom query
	if query != nil {
		q := req.URL.Query()
		for key, value := range query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Apply custom headers
	for header, value := range headers {
		req.Header.Set(header, value)
	}

	return req
}

func sendRequest(req *http.Request) (*http.Response, *errors.GenericError) {
	// Send request to API
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.WithField("error", err).Error("Failed to send HTTP request")
		logging.LogHTTPRequestResponse(req, res, log.ErrorLevel, "Send request failed")
		return nil, errors.NewGenericError(421, errorDomain, errorSubDomain, ErrorSendHTTPRequestFailed, nil)
	}

	// Dump request for debugging
	traceID, genErr := logging.LogHTTPRequestResponse(req, res, log.InfoLevel, "Debugging")
	if genErr != nil {
		log.WithField("error", genErr.GetDetailString()).Warn("Failed to log HTTP request and response")
	}

	// Handle Api errors
	switch {
	// API responded with an error
	case res.StatusCode >= 300:
		bodyBytes, genErr := duplicateAndReturnResponseBody(res, traceID)
		if genErr != nil {
			return nil, genErr
		}
		body := string(bodyBytes)

		// Log and return error
		log.WithFields(log.Fields{
			"body":     body,
			"trace_id": traceID,
		}).Warn("HTTP response code is error")
		meta := map[string]string{"response_body": body}
		return res, errors.NewGenericError(res.StatusCode, errorDomain, errorSubDomain, ErrorResponseCodeIsError, meta)

	// API responded with 2xx Success
	default:
		return res, nil
	}
}

// duplicateAndReturnResponseBody consumes the response body, replaces it with a new ReadCloser and returns the body
func duplicateAndReturnResponseBody(res *http.Response, traceID string) ([]byte, *errors.GenericError) {
	// Read response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// Reading body failed
		log.WithFields(log.Fields{
			"error":    err,
			"trace_id": traceID,
		}).Error("Unable to read HTTP response body")
		return nil, errors.NewGenericError(500, errorDomain, errorSubDomain, ErrorReadResponseBodyFailed, nil)
	}

	// Replace body with new reader
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Return body
	return body, nil
}

// UnwrapResponseCodeIsError checks if the error has code "response_code_is_error".
// If yes, it extracts the error message and builds the metadata for audit logging as well.
// If no, these will be filled with zero values.
func UnwrapResponseCodeIsError(ctx context.Context, genErr *errors.GenericError) (string, map[string]interface{}, bool) {
	// Return original error if not code "response_code_is_error"
	auditMeta := map[string]interface{}{}
	if genErr.SubDomainCode != ErrorResponseCodeIsError {
		return "", auditMeta, false
	}

	// Add response body to audit meta
	errorMsg := genErr.Meta["response_body"]
	auditMeta["response_body"] = errorMsg

	// Convert error message to lowercase for comparison
	return strings.ToLower(errorMsg), auditMeta, true
}
