package logging

import (
	"context"
	"strings"

	"github.com/micro/go-micro/v2/server"
	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/v2/converters"
	"github.com/skiprco/go-utils/v2/errors"
	"github.com/skiprco/go-utils/v2/metadata"
)

// AuditHandlerWrapper injects the service's name and called endpoint into the context.
// Next, it calls the required audit logging helpers.
//
// Usage:
//
// 	    service := micro.NewService(
//		    micro.Name(manifest.ServiceName),
//	    )
//	    service.Server().Init(server.WrapHandler(logging.AuditHandlerWrapper))
//	    service.Init()
func AuditHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		// Extract info from request
		service := strings.ToLower(req.Service())
		endpoint := strings.ToLower(req.Endpoint())

		// Inject info into context
		ctx, _, genErr := metadata.UpdateGoMicroMetadata(ctx, metadata.Metadata{
			"service_name":     service,
			"service_endpoint": endpoint,
		})
		if genErr != nil {
			return genErr.ToMicroError()
		}

		// Log attempt
		attempt := service + "_" + endpoint
		AuditAttempt(ctx, attempt, nil)

		// Call function
		err := fn(ctx, req, rsp)

		// Log result
		if err == nil {
			AuditSuccess(ctx, attempt, nil)
		} else {
			AuditFail(ctx, attempt, nil)
		}

		// Return result
		return err
	}
}

// AddAuditInfo prefixes the key with the service name, converts it to snake_case and adds the result to the context.
func AddAuditInfo(ctx context.Context, key string, value string) (context.Context, *errors.GenericError) {
	// Fetch service name
	serviceName, genErr := getFieldFromMeta(ctx, "service_name")
	if genErr != nil {
		return ctx, genErr
	}

	// Update metadata
	key = serviceName + "_" + key
	ctx, _, genErr = metadata.SetGoMicroMetadata(ctx, converters.ToSnakeCase(key), value)
	return ctx, genErr
}

// AddAuditInfoMap prefixes each metadata key with the service name, converts them to snake_case and adds the result to the context
func AddAuditInfoMap(ctx context.Context, info map[string]string) (context.Context, *errors.GenericError) {
	// Fetch service name
	serviceName, genErr := getFieldFromMeta(ctx, "service_name")
	if genErr != nil {
		return ctx, genErr
	}

	// Fix keys
	meta := make(metadata.Metadata, len(info))
	for k, v := range info {
		key := converters.ToSnakeCase(serviceName + "_" + k)
		meta[key] = v
	}

	// Update metadata
	ctx, _, genErr = metadata.UpdateGoMicroMetadata(ctx, meta)
	return ctx, genErr
}

func getFieldFromMeta(ctx context.Context, key string) (string, *errors.GenericError) {
	// Metadata from context
	meta, genErr := metadata.GetGoMicroMetadata(ctx)
	if genErr != nil {
		return "", genErr
	}

	// Extract service name
	serviceName := meta.Get(key)
	if serviceName == "" {
		log.WithField("meta", meta).WithField("key", key).Error("Key not found in context")
		return "", errors.NewGenericError(500, errDomain, errSubDomain, key+"_not_found_in_context", nil)
	}
	return serviceName, nil
}
