# go-utils
Common utils for Golang

## Usage
### Import
```go
import (
    "github.com/skiprco/go-utils/auth"
    "github.com/skiprco/go-utils/collections"
    "github.com/skiprco/go-utils/converters"
    "github.com/skiprco/go-utils/errors"
    "github.com/skiprco/go-utils/gin"
    "github.com/skiprco/go-utils/http"
    "github.com/skiprco/go-utils/logging"
    "github.com/skiprco/go-utils/manifest"
    "github.com/skiprco/go-utils/metadata"
    "github.com/skiprco/go-utils/validation"
)
```

### Auth
```go
// AuthOverridePrefix contains the prefix you have to use to override the authentication.
// This is needed when the action is not invoked by a user (e.g. callback by provider).
const AuthOverridePrefix = "system_override_"

// RoleUser is a user role which means the user is using the Skipr application
const RoleUser = "USER"

// RoleOperatorRead is a user role means the operator has read-only access to all data
const RoleOperatorRead = "OPERATOR_READ"

// RoleOperatorWrite is a user role means the operator has read-write access to all data
const RoleOperatorWrite = "OPERATOR_WRITE"

// RoleOperatorAdmin is a user role means the user has read-write access to all data and can modify the roles of other users
const RoleOperatorAdmin = "OPERATOR_ADMIN"

// HasRole checks if the provided roles is included in the user roles.
// This role might be granted implicitely (e.g. OPERATOR_READ on OPERATOR_ADMIN).
func HasRole(role string, userRoles []string) (bool, *errors.GenericError) {}

// IsOverride checks if the provided user ID has a prefix to override the authentication.
// Overrides will be clearly logged.
func IsOverride(ctx context.Context, userID string, subDomain string) bool {}

// HasRoleCheck is a helper function to validate if a user has a specific role
type HasRoleCheck func(ctx context.Context, userID string) (bool, *errors.GenericError)

// MustHaveRole is a helper to ease the implementation of access control checks.
// This helper should be called by a specific helper for the service which implements the access control.
func MustHaveRole(ctx context.Context, hasRoleCheck HasRoleCheck, userID string, errorDomain string, subDomain string) *errors.GenericError {}
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
code, genErr = converters.CountryNameToCountryCode("CURAÇAO") // code == "CW", genErr.Code is 404 when not found
```

#### Strings
```go
// Removes all the accents from the letters in the string, but keeps casing
normalised, genErr := converters.NormaliseString("Tëst Çôdé (test-result)") // normalised == "Test Code (test-result)"

// Converts a string to snake_case
output := converters.ToSnakeCase("ThisIS_a-veryRandom_string") // output == "this_is_a_very_random_string"

// Removes any character(included spaces) which is not a digit or a letter from a string
output := converters.CleanSpecialCharacters("dir.ty-Str*in//g :)") // output == "dirtyString"

// Sanitize removes all HTML tags from the input and escapes entities.
// Following entities are excluded from escaping: &, ' (apos)
output := converters.Sanitize("<p>Test</p>") // output == "Test"

// SanitizeObject takes a pointer to an object (struct, map, slice, ...) as input
// and runs Sanitize for each field which is a (pointer to a) string.
test := map[string]string{"<p>Key</p>": "<p>Value</p>"}
converters.SanitizeObject(&test) // test == map[string]string{"Key": "Value"}
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

// In case the server returns an error code (>= 300), the response body is present on
// the error meta as key "response_body". Also, the full response will still be returned.
// This way, you can translate the body to a more specific error using below setup.
res, genErr := http.Call("POST", "https://skipr.co", "/test", request, response, nil, nil)
if genErr != nil {
    // UnwrapResponseCodeIsError checks if the error has code "response_code_is_error".
    // If yes, it extracts the error message and builds the metadata for audit logging as well.
    // If no, these will be filled with zero values.
    //
    // Check if genErr has subdomain code "response_code_is_error"
    errorMsg, auditMeta, ok := http.UnwrapResponseCodeIsError(ctx, genErr)
    logging.AuditFail(ctx, "...", auditMeta)
    if !ok {
        return nil, genErr
    }

    // Translate error to a more specific one
    switch {
        case strings.Contains(errorMsg, "..."):
            ...
    }

    // No specific error matched => Return original
    return nil, genErr
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

// Add the AuditHandlerWrapper to a service
service := micro.NewService(
    micro.Name(manifest.ServiceName),
)
service.Server().Init(server.WrapHandler(logging.AuditHandlerWrapper))
service.Init()

// AddAuditInfo prefixes the key with the service name, converts it to snake_case and adds the result to the context.
func AddAuditInfo(ctx context.Context, key string, value string) (context.Context, *errors.GenericError) {}

// AddAuditInfoMap prefixes each metadata key with the service name, converts them to snake_case and adds the result to the context
func AddAuditInfoMap(ctx context.Context, info map[string]string) (context.Context, *errors.GenericError) {}
```

### Manifest
```go
// Load manifest file
manifest, genErr := manifest.LoadManifest()
if genErr != nil {
    log.WithField("error", genErr.GetDetailString()).Panic("Failed to load manifest file")
}
```

### Metadata
```go
// Get tries to fetch the value stored with the provided key.
// If meta is nil or key doesn't exist, an empty string is returned.
func (meta Metadata) Get(key string) string {}

// GetGinMetadata returns the currently defined metadata from the Gin context
func GetGinMetadata(c *gin.Context) Metadata {}

// UpdateGinMetadata upserts the metadata stored in the Gin context. Returns result of the merge.
func UpdateGinMetadata(c *gin.Context, additionalMetadata Metadata) Metadata {}

// SetGinMetadata upserts a single key/value pair in the Gin context. Returns result of the merge.
func SetGinMetadata(c *gin.Context, key string, value string) Metadata {}

// GetGoMicroMetadata returns the currently defined metadata from the go-micro context
func GetGoMicroMetadata(ctx context.Context) (Metadata, *errors.GenericError) {}

// GetUserIDFromGoMicroMeta extracts the user ID from the metadata of go-micro
// and also returns the raw metadata for later use. Throws an error if unable
// to read metadata or if user_id is not set.
func GetUserIDFromGoMicroMeta(ctx context.Context, errorDomain string) (string, Metadata, *errors.GenericError) {}

// UpdateGoMicroMetadata upserts the metadata stored in the go-micro context. Returns result of the merge.
func UpdateGoMicroMetadata(ctx context.Context, additionalMetadata Metadata) (context.Context, Metadata, *errors.GenericError) {}

// SetGoMicroMetadata upserts a single key/value pair in the go-micro context. Returns result of the merge.
func SetGoMicroMetadata(ctx context.Context, key string, value string) (context.Context, Metadata, *errors.GenericError) {}

// ConvertGinToGoMicro returns a context object with the metadata
// set as go-micro metadata. This way the metadata can be accessed in
// each microservices.
func ConvertGinToGoMicro(c *gin.Context) context.Context {}
```

### Validation

#### Country code
```go
// Checks if a country code is valid. An empty code is considered valid as well.
valid := validation.ValidateCountryCode("BE") // valid == true
```

#### Country code
```go
// ValidateAndFormatPhoneNumber checks if the provided phone number is valid.
// If yes, it will format it to its E164 representation (e.g +32...)
number, genErr := validation.ValidateAndFormatPhoneNumber("+32 478 12 34 56") // number == "+32478123456"
```

#### Time
```go
// WithinTimeRange checks if the "nowTime" lies between "startTime" (including) and "endTime" (excluding).
now := time.Now()
start := now.Add(-time.Hour)
end := now.Add(time.Hour)
within, genErr := validation.WithinTimeRange(now, start, end) // within = true
```
### mongo
```go
// init mongo client & collection
repo, genErr := repository.NewMongoRepository(context.Background(), "mongoAddress", "mongoDBName", []string{"CollectionName"})

// query
query := map[string]interface{}{
		"id":      "123456",
		"status": "pending",
	}

// count
count := repo.Count(ctx, "CollectionName", query, "functionName")

// get one function call
result := &myEntityStruct{}
genErr := repo.GetOne(ctx, "CollectionName", query, true, result,"functionName")
if genErr != nil {
    // manage genErr
    
}
// result is populate from database information
if result.Id != "" {
    // Be careful : as acceptsEmptyResult parameter is true, if the entity is not found the field of result will empty
}

// get multiple function call
 query := map[string]interface{}{
 		"status": "pending",
 	}
result := &[]myEntityStruct{}
genErr := repo.GetMultiple(ctx, "CollectionName", query, results,"functionName")

// Save & delete
repo.Save(ctx, "CollectionName", myEntity, myEntity.Id, "functionName")
repo.Delete(ctx, "CollectionName", myEntity.Id, "functionName")
```