# go-utils
Common utils for Golang

## Usage
### Import
```go
import (
    "github.com/skiprco/go-utils/errors"
    "github.com/skiprco/go-utils/http"
    "github.com/skiprco/go-utils/logging"
    "github.com/skiprco/go-utils/manifest"
)
```

### Errors
```go
// Create an error with metadata
meta := map[string]string{
    "provider": "go-utils",
}
genErr := errors.NewGenericError(500, "booking", "common", "marshal_request_body_failed", meta)

// Create an error from a go-micro error
genErr := errors.NewGenericFromMicroError(microError)
```

### HTTP
```go
// Call with JSON marshaling and parsing
users := []User{...}
response := CreateUserResponse{}
query = map[string]string{"lang": "en"}
headers = map[string]string{"Authorization": "Bearer ..."}
genErr := http.Call("POST", "https://skipr.co", "users", users, CreateUserResponse, query, headers)

// Call with raw body and response (bytes)
file = []byte{...}
httpResponse, genErr := http.CallRaw("POST", "https://skipr.co", "files", file, nil, nil)

// In case the server returns an error code (>= 300), the http.Response is still returned.
// This way, you can translate the body to a more specific error using below setup
res, genErr := Call("POST", "https://skipr.co", "/test", request, response, nil, nil)
if genErr != nil {
    if genErr.SubDomainCode != "response_code_is_error" {
        return nil, genErr
    }
    return nil, translateResponseToSpecificError(res)
}
```

### Logging
```go
import log "github.com/sirupsen/logrus"

// Setup default Logrus logger
//   - Log in JSON format
//   - Split Error logs to stderr, others to stdout
//   - Include stacktrace on log level Error
logging.SetupLogger("replace_me_with_service_name")

// Log a generic error
err := ...
log.WithField("error", err).Error("Human readable message")

// Log HTTP request and response
req, res := ...
logging.LogHTTPRequestResponse(req, res, log.WarnLevel, "Human readable warning message")
logging.LogHTTPRequestResponse(req, res, log.ErrorLevel, "Human readable error message")
```

### Manifest
```go
// Load manifest file
manifest, genErr := manifest.LoadManifest()
if genErr != nil {
    log.WithField("error", genErr.GetDetailString()).Panic("Failed to load manifest file")
}
```
