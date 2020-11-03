package logging

import (
	"context"
	"strings"

	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
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
		ctx = metadata.Set(ctx, "service_name", service)
		ctx = metadata.Set(ctx, "service_endpoint", endpoint)

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
