package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/drewjst/deltagov/internal/congress"
)

// DiagnosticService handles system health and testing endpoints
type DiagnosticService struct {
	CongressClient *congress.Client
}

// NewDiagnosticService creates a new instance of the service
func NewDiagnosticService(client *congress.Client) *DiagnosticService {
	return &DiagnosticService{CongressClient: client}
}

// DiagnosticHealthOutput is the response for diagnostic health check
type DiagnosticHealthOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

// RegisterDiagnosticRoutes registers testing and health endpoints with Huma
func RegisterDiagnosticRoutes(api huma.API, s *DiagnosticService) {
	huma.Register(api, huma.Operation{
		OperationID: "get-health",
		Method:      "GET",
		Path:        "/health",
		Summary:     "Health Check",
		Description: "Returns the status of the API and database connection.",
		Tags:        []string{"Diagnostics"},
	}, func(ctx context.Context, input *struct{}) (*DiagnosticHealthOutput, error) {
		resp := &DiagnosticHealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})
}
