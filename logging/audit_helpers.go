package logging

// See https://www.notion.so/skipr/Logging-Technical-Doc-b12c01973e3046daa82f98b51fa06251 for more info

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/metadata"
	log "github.com/sirupsen/logrus"
)

// ========================================
// =               CONSTANTS              =
// ========================================

type auditCategory string

const (
	auditCategoryFact auditCategory = "fact"

	auditCategoryAttempt auditCategory = "attempt"
	auditCategorySuccess auditCategory = "success"
	auditCategoryFail    auditCategory = "fail"
)

// ========================================
// =         START / END OPERATION        =
// ========================================

// AuditOperationAttempt logs an audit message in the "attempt" category for an operation.
// This call should be used at the start of an operation.
func AuditOperationAttempt(ctx context.Context, additionalData map[string]interface{}) {
	AuditAttempt(ctx, "operation_attempt", additionalData)
}

// AuditOperationSuccess logs a successful result for a previous call to AuditOperationAttempt.
// This call should be used at the end of an operation.
func AuditOperationSuccess(ctx context.Context, additionalData map[string]interface{}) {
	AuditSuccess(ctx, "operation_success", additionalData)
}

// AuditOperationFail logs a failure result for a previous call to AuditOperationAttempt
// This call should be used at the end of an operation.
func AuditOperationFail(ctx context.Context, additionalData map[string]interface{}) {
	AuditFail(ctx, "operation_fail", additionalData)
}

// ========================================
// =           DURING OPERATION           =
// ========================================

// AuditFact logs an audit message in the "fact" category
func AuditFact(ctx context.Context, message string, additionalData map[string]interface{}) {
	logEvent(ctx, message, auditCategoryFact, additionalData)
}

// AuditAttempt logs an audit message in the "attempt" category
func AuditAttempt(ctx context.Context, attemptName string, additionalData map[string]interface{}) {
	logEvent(ctx, attemptName, auditCategoryAttempt, additionalData)
}

// AuditSuccess logs a successful result to a previous audit message of category "attempt"
func AuditSuccess(ctx context.Context, attemptName string, additionalData map[string]interface{}) {
	logEvent(ctx, attemptName, auditCategorySuccess, additionalData)
}

// AuditFail logs a failure result to a previous audit message of category "attempt"
func AuditFail(ctx context.Context, attemptName string, additionalData map[string]interface{}) {
	logEvent(ctx, attemptName, auditCategoryFail, additionalData)
}

// ========================================
// =                HELPERS               =
// ========================================

func logEvent(ctx context.Context, message string, category auditCategory, additionalData map[string]interface{}) {
	// Log priority
	// A lower priority (e.g. 3) will be overwritten by higher priority (e.g. 1)
	//
	// 1. Directly provided data: message and category
	// 2. Derived data: operation_time
	// 3. Additionally provided data: additionalData
	// 4. Metadata present in context: ctx

	// Read metadata from context
	var logFields map[string]interface{}
	meta, success := metadata.FromContext(ctx)
	if success && meta != nil {
		maxLength := len(meta) + len(additionalData) + 2 // for message and category
		logFields = make(log.Fields, maxLength)
		for key, value := range meta {
			lowerKey := strings.ToLower(key)
			logFields[lowerKey] = value
		}
	} else {
		maxLength := len(additionalData) + 2 // for message and category
		logFields = make(log.Fields, maxLength)
	}

	// Add additional data
	for key, value := range additionalData {
		logFields[key] = value
	}

	// Derive data
	deriveOperationTime(logFields)

	// Add directly provided data
	logFields["message"] = message
	logFields["category"] = category

	// Send audit to log
	log.WithFields(logFields).Info(message)
}

// Calculate operation_time based on operation_start_datetime and set in fields
func deriveOperationTime(fields log.Fields) {
	// Extract start from fields
	startInterface, ok := fields["operation_start_datetime"]
	if !ok {
		return
	}

	// Prepare logging
	deriveLog := log.WithFields(log.Fields{
		"field": "operation_start_datetime",
		"value": startInterface,
	})

	// Convert value to string
	startString, ok := startInterface.(string)
	if !ok {
		deriveLog.WithField("type", reflect.TypeOf(startInterface)).Error("Expected type string for field operation_start_datetime")
		return
	}

	// Parse start to time
	start, err := time.Parse(time.RFC3339, startString)
	if err != nil {
		deriveLog.Error("Unable to parse operation_start_datetime as RFC3339")
		return
	}

	// Time should not be initial
	if start.IsZero() {
		deriveLog.Error("operation_start_datetime should not contain a zero value")
		return
	}

	// Derive and set operation_time
	operationTime := time.Now().Sub(start).Milliseconds()
	fields["operation_time"] = strconv.FormatInt(operationTime, 10)
}
