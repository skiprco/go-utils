# go-utils
Common utils for Golang

## Usage
### Import
```go
import (
    "github.com/skiprco/go-utils/errors"
    "github.com/skiprco/go-utils/logging"
)
```

### Errors

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
