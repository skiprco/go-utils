package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/v2/errors"
	"github.com/skiprco/go-utils/v2/logging"
)

// Call marshals the body to JSON and sends a HTTP request. Response is parsed as JSON.
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
			return nil, errors.NewGenericError(500, "go_utils", "common", "marshal_request_body_failed", nil)
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
		return nil, errors.NewGenericError(421, "go_utils", "common", "read_response_body_failed", nil)
	}

	// Unmarshal response
	err = json.Unmarshal(resBody, response)
	if err != nil {
		traceLog.WithField("error", err).Error("Failed to parse response body from JSON")
		return nil, errors.NewGenericError(421, "go_utils", "common", "parse_response_body_failed", nil)
	}

	// Return result
	return res, nil
}

// CallRaw sends a HTTP request without modifying the body
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
		logging.LogHTTPRequestResponse(req, res, log.ErrorLevel, "Send request failed")
		return nil, errors.NewGenericError(421, "go_utils", "common", "send_http_request_failed", nil)
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
		body, genErr := duplicateAndReturnResponseBody(res, traceID)
		if genErr != nil {
			return nil, genErr
		}

		// Log and return error
		log.WithFields(log.Fields{
			"body":     string(body),
			"trace_id": traceID,
		}).Warn("HTTP response code is error")
		return res, errors.NewGenericError(res.StatusCode, "go_utils", "common", "response_code_is_error", nil)

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
		return nil, errors.NewGenericError(500, "go_utils", "common", "read_response_body_failed", nil)
	}

	// Replace body with new reader
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Return body
	return body, nil
}
