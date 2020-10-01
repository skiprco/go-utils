package logging

import (
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/errors"
)

// LogHTTPRequestResponse logs a HTTP request and response
func LogHTTPRequestResponse(request *http.Request, response *http.Response, logLevel log.Level, message string) *errors.GenericError {
	// Dump request
	dumpReq, err := httputil.DumpRequest(request, true)
	if err != nil {
		log.WithField("error", err).Error("Unable to dump HTTP request")
		return errors.NewGenericError(500, "go_utils", "logging", "unable_to_dump_request", nil)
	}

	// Response body is already consumed.
	// Therefore, consuming the body will always fail
	// => Changed "body" parameter of DumpResponse to false
	var dumpResponse string
	statusCode := 0
	if response != nil {
		statusCode = response.StatusCode
		dumpResponseByte, err := httputil.DumpResponse(response, false)
		if err != nil {
			log.WithField("error", err).Error("Unable to dump HTTP response")
			return errors.NewGenericError(500, "go_utils", "logging", "unable_to_dump_response", nil)
		}
		dumpResponse = string(dumpResponseByte)
	} else {
		dumpResponse = "response_is_nil"
	}

	// Log result
	log.WithFields(log.Fields{
		"request":            string(dumpReq),
		"response":           dumpResponse,
		"response_http_code": statusCode,
	}).Log(logLevel, message)

	// Logging successful
	return nil
}
