package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/skiprco/go-utils/errors"
	"github.com/stretchr/testify/assert"
)

// ########################################
// #                COMMON                #
// ########################################
type requestSample struct {
	Request string `json:"request"`
}

type responseSample struct {
	Response string `json:"response"`
}

type TestHandlerFunc func(t *testing.T, res http.ResponseWriter, req *http.Request) bool

func cleanJSON(s string) string {
	re := regexp.MustCompile(`[^\S ]`)
	return re.ReplaceAllString(s, "")
}

func getAPIServerMock(t *testing.T, testHandlerFunc TestHandlerFunc) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		fmt.Println(req.Method, req.URL.Path)

		// Call test specific handler
		if !testHandlerFunc(t, res, req) {
			assert.Failf(t, "Unhandled call to test server",
				"Code called \"%s %s\", which is not handled by the test server", req.Method, req.URL.Path)
		}
	}))
	return ts
}

// ########################################
// #                 Call                 #
// ########################################
func Test_Call_200(t *testing.T) {
	// Setup mock handlers
	testFunc := func(t *testing.T, res http.ResponseWriter, req *http.Request) bool {
		switch req.URL.Path {
		case "/test":
			assert.Equal(t, "TESTMETHOD", req.Method)
			assert.Equal(t, "query-success", req.URL.Query().Get("test-query"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, "header-success", req.Header.Get("X-TEST"))
			fmt.Fprint(res, `{"Response": "response sample"}`)
			return true
		}

		// Call not handled
		return false
	}

	// Create mock server
	ts := getAPIServerMock(t, testFunc)
	defer ts.Close()

	// Setup request and response
	request := &requestSample{
		Request: "request sample",
	}
	response := &responseSample{}
	query := map[string]string{"test-query": "query-success"}
	headers := map[string]string{"X-TEST": "header-success"}

	// Call helper
	resp, genErr := Call("TESTMETHOD", ts.URL, "/test", request, response, query, headers)

	// Assert results
	expectedResponse := &responseSample{
		Response: "response sample",
	}
	assert.Nil(t, genErr)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedResponse, response)
}

func Test_Call_400(t *testing.T) {
	// Setup mock handlers
	testFunc := func(t *testing.T, res http.ResponseWriter, req *http.Request) bool {
		switch req.URL.Path {
		case "/test":
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			res.WriteHeader(400)
			fmt.Fprint(res, "error_during_test")
			return true
		}

		// Call not handled
		return false
	}

	// Create mock server
	ts := getAPIServerMock(t, testFunc)
	defer ts.Close()

	// Setup request and response
	response := &responseSample{}

	// Call helper
	resp, genErr := Call("POST", ts.URL, "/test", nil, response, nil, nil)

	// Assert results
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "error_during_test", string(body))
	assert.Equal(t, &responseSample{}, response)
	errors.AssertGenericError(t, genErr, 400, "response_code_is_error", nil)
}

func Test_Call_MarshalRequestBodyFailed_Failure(t *testing.T) {
	// Call helper
	request := func() {}
	response := &responseSample{}
	resp, genErr := Call("POST", "https://skipr.co", "/test", request, response, nil, nil)

	// Assert results
	assert.Nil(t, resp)
	assert.Equal(t, &responseSample{}, response)
	errors.AssertGenericError(t, genErr, 500, "marshal_request_body_failed", nil)
}

func Test_Call_SendRequestFailed_Failure(t *testing.T) {
	// Call helper
	response := &responseSample{}
	resp, genErr := Call("POST", "invalid://skipr.co", "/test", nil, response, nil, nil)

	// Assert results
	assert.Nil(t, resp)
	assert.Equal(t, &responseSample{}, response)
	errors.AssertGenericError(t, genErr, 421, "send_http_request_failed", nil)
}

func Test_Call_ParseResponseBodyFailed_Failure(t *testing.T) {
	// Setup mock handlers
	testFunc := func(t *testing.T, res http.ResponseWriter, req *http.Request) bool {
		switch req.URL.Path {
		case "/test":
			fmt.Fprint(res, "invalid")
			return true
		}

		// Call not handled
		return false
	}

	// Create mock server
	ts := getAPIServerMock(t, testFunc)
	defer ts.Close()

	// Call helper
	response := &responseSample{}
	resp, genErr := Call("POST", ts.URL, "/test", nil, response, nil, nil)

	// Assert results
	assert.Nil(t, resp)
	assert.Equal(t, &responseSample{}, response)
	errors.AssertGenericError(t, genErr, 421, "parse_response_body_failed", nil)
}

// ########################################
// #                CallRaw               #
// ########################################
func Test_CallRaw_200(t *testing.T) {
	// Setup mock handlers
	testFunc := func(t *testing.T, res http.ResponseWriter, req *http.Request) bool {
		switch req.URL.Path {
		case "/test":
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.Equal(t, body, []byte{1, 2, 3})
			assert.Equal(t, "TESTMETHOD", req.Method)
			assert.Equal(t, "query-success", req.URL.Query().Get("test-query"))
			assert.Equal(t, "header-success", req.Header.Get("X-TEST"))
			return true
		}

		// Call not handled
		return false
	}

	// Create mock server
	ts := getAPIServerMock(t, testFunc)
	defer ts.Close()

	// Setup request
	request := []byte{1, 2, 3}
	query := map[string]string{"test-query": "query-success"}
	headers := map[string]string{"X-TEST": "header-success"}

	// Call helper
	resp, genErr := CallRaw("TESTMETHOD", ts.URL, "/test", request, query, headers)

	// Assert results
	assert.Nil(t, genErr)
	assert.NotNil(t, resp)
}

func Test_CallRaw_400(t *testing.T) {
	// Setup mock handlers
	testFunc := func(t *testing.T, res http.ResponseWriter, req *http.Request) bool {
		switch req.URL.Path {
		case "/test":
			res.WriteHeader(400)
			fmt.Fprint(res, "error_during_test")
			return true
		}

		// Call not handled
		return false
	}

	// Create mock server
	ts := getAPIServerMock(t, testFunc)
	defer ts.Close()

	// Call helper
	resp, genErr := CallRaw("GET", ts.URL, "/test", nil, nil, nil)

	// Assert results
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "error_during_test", string(body))
	errors.AssertGenericError(t, genErr, 400, "response_code_is_error", nil)
}
