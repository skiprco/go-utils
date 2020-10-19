package logging

import (
	"io/ioutil"
	"os"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
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

	// Add error hook to format log for error reporting
	log.AddHook(&TraceOnErrorHook{service: serviceName})

	// Split output between stderr and stdout
	log.SetOutput(ioutil.Discard) // Send all logs to nowhere by default
	log.AddHook(&writer.Hook{     // Send logs with level higher than error to stderr
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		},
	})
	log.AddHook(&writer.Hook{ // Send warning and lower logs to stdout
		Writer: os.Stdout,
		LogLevels: []log.Level{
			log.WarnLevel,
			log.InfoLevel,
			log.DebugLevel,
		},
	})

	// Only log the info severity or above
	log.SetLevel(log.InfoLevel)
}

// TraceOnErrorHook will add a stack trace when a error is logged
type TraceOnErrorHook struct {
	service string
}

// Levels returns the levels on which TraceOnErrorHook should be triggered
func (h *TraceOnErrorHook) Levels() []log.Level {
	return []log.Level{
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}

// Fire adds the stacktrace to the log entry when triggered
func (h *TraceOnErrorHook) Fire(entry *log.Entry) error {
	entry.Data["serviceName"] = h.service
	entry.Data["stack_trace"] = string(debug.Stack())
	entry.Data["@type"] = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
	return nil
}
