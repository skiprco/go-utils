# go-utils
Common utils for Golang

## Usage
### Import
```go
import (
    "github.com/skiprco/go-utils/collections"
    "github.com/skiprco/go-utils/converters"
    "github.com/skiprco/go-utils/errors"
    "github.com/skiprco/go-utils/gin"
    "github.com/skiprco/go-utils/http"
    "github.com/skiprco/go-utils/logging"
    "github.com/skiprco/go-utils/manifest"
    "github.com/skiprco/go-utils/validation"
)
```

### Collections

#### String Map: map[string]string
```go
// StringMapMerge: Merge 2 maps into a copy
base := map[string]string{"k1": "b1", "k2": "b2"}
additional := map[string]string{"k2": "a2", "k3": "a3"}
result := collections.StringMapMerge(base, additional) // == map[string]string{"k1": "b1", "k2": "a2", "k3": "a3"}
```

#### String Slice: []string
```go
// StringSliceContains: Check if the slice contains a value
slice := []string{"test1", "test2"}
result := collections.StringSliceContains(slice, "test1") // == true
```

### Converters

#### Country code
```go
// Fetch map country code to country name
cc := converters.CountryCodes()
name := cc["BE"] // name = "Belgium", nil if not found

// Convert country code to country name
name, genErr := converters.CountryCodeToCountryName("BE") // name == "Belgium", genErr.Code is 404 when not found

// Convert a country's name to a country's code, ignoring casing and accents
code, genErr = CountryNameToCountryCode("CURAÇAO") // code == "CW", genErr.Code is 404 when not found
```

#### Strings
```go
// Removes all the accents from the letters in the string, but keeps casing
normalised, genErr := NormaliseString("Tëst Çôdé (test-result)") // normalised == "Test Code (test-result)"
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

### Gin
```go
// Global audit middleware for Gin
// - Logs each request and response
// - Injects metadata into the context to support audit logging in other services
router.Use(gin.AuditMiddleware("booking-api"))

// Fetch current metadata from Gin context
meta := GetMetadata(c)

// Update metadata in Gin context
additionalMeta := map[string]string{ ... }
updatedMeta := UpdateMetadata(c, additionalMeta)

// Fetch context with metadata from Gin context
ctx := GetContextWithMetadata(c)
response, err := or.OrganisationClient.CreateOrganisationFromVAT(ctx, request)
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

// Create a new audit log entry
additionalData := map[string]interface{} {
    "reason": "example reason",
}
logging.AuditFact(ctx, "user_update_skipped", additionalData)
logging.AuditAttempt(ctx, "update_user", nil)
logging.AuditSuccess(ctx, "update_user", nil)
logging.AuditFail(ctx, "update_user", nil)
```

### Manifest
```go
// Load manifest file
manifest, genErr := manifest.LoadManifest()
if genErr != nil {
    log.WithField("error", genErr.GetDetailString()).Panic("Failed to load manifest file")
}
```

### Validation

#### Country code
```go
// Checks if a country code is valid. An empty code is considered valid as well.
valid := ValidateCountryCode("BE") // valid == true
```
