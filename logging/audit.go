package logging

// See https://www.notion.so/skipr/Logging-Technical-Doc-b12c01973e3046daa82f98b51fa06251 for more info

import (
	"context"
	"strings"

	"github.com/micro/go-micro/v2/metadata"
	log "github.com/sirupsen/logrus"
)

type auditCategory string

const (
	auditCategoryFact auditCategory = "fact"

	auditCategoryAttempt auditCategory = "attempt"
	auditCategorySuccess auditCategory = "success"
	auditCategoryFail    auditCategory = "fail"
)

// AuditFact logs an audit message in the "fact" category
func AuditFact(ctx context.Context, message string, additionalData map[string]interface{}) {
	logEvent(ctx, message, auditCategoryFact, additionalData)
}

// AuditAttempt logs an audit message in the "attempt" category
func AuditAttempt(ctx context.Context, message string, additionalData map[string]interface{}) {
	logEvent(ctx, message, auditCategoryAttempt, additionalData)
}

// AuditSuccess logs a successful result to a previous audit message of category "attempt"
func AuditSuccess(ctx context.Context, message string, additionalData map[string]interface{}) {
	logEvent(ctx, message, auditCategorySuccess, additionalData)
}

// AuditFail logs a failure result to a previous audit message of category "attempt"
func AuditFail(ctx context.Context, message string, additionalData map[string]interface{}) {
	logEvent(ctx, message, auditCategoryFail, additionalData)
}

func logEvent(ctx context.Context, message string, category auditCategory, additionalData map[string]interface{}) {
	// Log priority
	// A lower priority (e.g. 3) will be overwritten by higher priority (e.g. 1)
	//
	// 1. Directly provided data: message and category
	// 2. Additionally provided data: additionalData
	// 3. Metadata present in context: ctx

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

	// Add directly provided data
	logFields["message"] = message
	logFields["category"] = category

	// Send audit to log
	log.WithFields(logFields).Info(message)
}
