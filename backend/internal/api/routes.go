package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// --- Request/Response Types ---

// ListBillsOutput is the response for listing bills
type ListBillsOutput struct {
	Body struct {
		Bills []MockBill `json:"bills"`
		Total int        `json:"total"`
	}
}

// GetBillInput is the request for getting a single bill
type GetBillInput struct {
	ID string `path:"id" doc:"Bill ID (e.g., hr1234-119)"`
}

// GetBillOutput is the response for getting a single bill
type GetBillOutput struct {
	Body MockBill
}

// GetBillVersionsInput is the request for getting bill versions
type GetBillVersionsInput struct {
	ID string `path:"id" doc:"Bill ID"`
}

// GetBillVersionsOutput is the response for getting bill versions
type GetBillVersionsOutput struct {
	Body struct {
		BillID   string        `json:"billId"`
		Versions []MockVersion `json:"versions"`
	}
}

// ComputeDiffInput is the request for computing a diff
type ComputeDiffInput struct {
	BillID      string `path:"billId" doc:"Bill ID"`
	FromVersion string `path:"fromVersion" doc:"Source version ID"`
	ToVersion   string `path:"toVersion" doc:"Target version ID"`
}

// ComputeDiffOutput is the response for computing a diff
type ComputeDiffOutput struct {
	Body MockDelta
}

// HealthOutput is the response for health check
type HealthOutput struct {
	Body struct {
		Status  string `json:"status"`
		Service string `json:"service"`
	}
}

// --- Route Registration ---

// RegisterRoutes sets up all API routes with Huma
func RegisterRoutes(api huma.API) {
	// Health check
	huma.Get(api, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "healthy"
		resp.Body.Service = "deltagov-api"
		return resp, nil
	})

	// List all bills
	huma.Get(api, "/api/v1/bills", func(ctx context.Context, input *struct{}) (*ListBillsOutput, error) {
		bills := GetMockBills()
		resp := &ListBillsOutput{}
		resp.Body.Bills = bills
		resp.Body.Total = len(bills)
		return resp, nil
	})

	// Get single bill
	huma.Register(api, huma.Operation{
		OperationID: "get-bill",
		Method:      http.MethodGet,
		Path:        "/api/v1/bills/{id}",
		Summary:     "Get a bill by ID",
		Description: "Returns detailed information about a specific legislative bill",
		Tags:        []string{"Bills"},
	}, func(ctx context.Context, input *GetBillInput) (*GetBillOutput, error) {
		bills := GetMockBills()
		for _, bill := range bills {
			if bill.ID == input.ID {
				return &GetBillOutput{Body: bill}, nil
			}
		}
		return nil, huma.Error404NotFound("bill not found")
	})

	// Get bill versions
	huma.Register(api, huma.Operation{
		OperationID: "get-bill-versions",
		Method:      http.MethodGet,
		Path:        "/api/v1/bills/{id}/versions",
		Summary:     "Get all versions of a bill",
		Description: "Returns all tracked versions/snapshots of a bill's text",
		Tags:        []string{"Bills"},
	}, func(ctx context.Context, input *GetBillVersionsInput) (*GetBillVersionsOutput, error) {
		versions := GetMockVersions(input.ID)
		resp := &GetBillVersionsOutput{}
		resp.Body.BillID = input.ID
		resp.Body.Versions = versions
		return resp, nil
	})

	// Compute diff between versions
	huma.Register(api, huma.Operation{
		OperationID: "compute-diff",
		Method:      http.MethodGet,
		Path:        "/api/v1/bills/{billId}/diff/{fromVersion}/{toVersion}",
		Summary:     "Compute diff between two bill versions",
		Description: "Returns a structured diff showing insertions, deletions, and unchanged text between two versions",
		Tags:        []string{"Diff"},
	}, func(ctx context.Context, input *ComputeDiffInput) (*ComputeDiffOutput, error) {
		delta := GetMockDelta(input.FromVersion, input.ToVersion)
		return &ComputeDiffOutput{Body: delta}, nil
	})
}
