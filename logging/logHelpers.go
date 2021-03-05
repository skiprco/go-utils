package logging

import (
	"net/http"
	"net/http/httputil"

	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/v2/errors"
)

// LogHTTPRequestResponse logs a HTTP request and response. Returns a trace ID to link follow-up logs.
//
// Raises
//
// - 500/unable_to_dump_request: Failed to dump to HTTP request
//
// - 500/unable_to_dump_response: Failed to dump to HTTP response
func LogHTTPRequestResponse(request *http.Request, response *http.Response, logLevel log.Level, message string) (string, *errors.GenericError) {
	// Generate trace ID
	traceID := uuid.New()
	traceLog := log.WithField("trace_id", traceID)

	// Dump request
	dumpReq, err := httputil.DumpRequest(request, true)
	if err != nil {
		traceLog.WithField("error", err).Error("Unable to dump HTTP request")
		return traceID, errors.NewGenericError(500, errorDomain, errorSubDomain, ErrorUnableToDumpRequest, nil)
	}

	// Response body is already consumed.
	// Therefore, consuming the body will always fail
	// => Changed "body" parameter of DumpResponse to false
	var dumpResponse string
	statusCode := 0
	if response != nil {
		statusCode = response.StatusCode
		dumpResponseByte, err := httputil.DumpResponse(response, true)
		if err != nil {
			traceLog.WithField("error", err).Error("Unable to dump HTTP response")
			return traceID, errors.NewGenericError(500, errorDomain, errorSubDomain, ErrorUnableToDumpResponse, nil)
		}
		dumpResponse = string(dumpResponseByte)
	} else {
		dumpResponse = "response_is_nil"
	}

	// Log result
	traceLog.WithFields(log.Fields{
		"request":            string(dumpReq),
		"response":           dumpResponse,
		"response_http_code": statusCode,
		"context": map[string]interface{}{ // GCP specific fields
			"httpRequest": map[string]interface{}{
				"method":             request.Method,
				"url":                request.URL.String(),
				"responseStatusCode": statusCode,
			},
		},
	}).Log(logLevel, message)

	// Logging successful
	return traceID, nil
}
