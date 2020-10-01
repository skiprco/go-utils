package logging

import (
	"bytes"
	"os"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

// SetupLogger configures the standard logger of Logrus
func SetupLogger(serviceName string) {
	// Log as JSON instead of the default ASCII formatter
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "time",
			log.FieldKeyLevel: "severity",
			log.FieldKeyMsg:   "message",
			log.FieldKeyFunc:  "@caller",
		},
	})

	// set output to stderr or stdout
	log.SetOutput(&OutputSplitter{})

	// Only log the info severity or above
	log.SetLevel(log.InfoLevel)

	// add error hook to format log for error reporting
	log.AddHook(&TraceOnErrorHook{service: serviceName})
}

// OutputSplitter routes log output to stdout or stderr based on the level
type OutputSplitter struct{}

func (splitter *OutputSplitter) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte(`"severity":"error"`)) {
		return os.Stderr.Write(p)
	}
	return os.Stdout.Write(p)
}

// TraceOnErrorHook will add a stack trace when a error is logged
type TraceOnErrorHook struct {
	service string
}

// Levels returns the levels on which TraceOnErrorHook should be triggered
func (h *TraceOnErrorHook) Levels() []log.Level {
	return []log.Level{
		log.ErrorLevel,
	}
}

// Fire adds the stacktrace to the log entry when triggered
func (h *TraceOnErrorHook) Fire(entry *log.Entry) error {
	entry.Data["serviceName"] = h.service
	entry.Data["stack_trace"] = string(debug.Stack())
	entry.Data["@type"] = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
	return nil
}
