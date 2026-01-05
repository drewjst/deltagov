package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
)

// --- Request/Response Types ---

// ListBillsOutput is the response for listing bills
type ListBillsOutput struct {
	Body struct {
		Bills []BillResponse `json:"bills"`
		Total int            `json:"total"`
	}
}

// GetBillInput is the request for getting a single bill
type GetBillInput struct {
	ID uint `path:"id" doc:"Bill ID (database ID)"`
}

// GetBillOutput is the response for getting a single bill
type GetBillOutput struct {
	Body BillResponse
}

// GetBillVersionsInput is the request for getting bill versions
type GetBillVersionsInput struct {
	ID uint `path:"id" doc:"Bill ID"`
}

// GetBillVersionsOutput is the response for getting bill versions
type GetBillVersionsOutput struct {
	Body struct {
		BillID   uint              `json:"billId"`
		Versions []VersionResponse `json:"versions"`
	}
}

// ComputeDiffInput is the request for computing a diff
type ComputeDiffInput struct {
	BillID      uint `path:"billId" doc:"Bill ID"`
	FromVersion uint `path:"fromVersion" doc:"Source version ID"`
	ToVersion   uint `path:"toVersion" doc:"Target version ID"`
}

// ComputeDiffOutput is the response for computing a diff
type ComputeDiffOutput struct {
	Body DiffResponse
}

// HealthOutput is the response for health check
type HealthOutput struct {
	Body struct {
		Status  string `json:"status"`
		Service string `json:"service"`
	}
}

// FetchHR1Output is the response for fetching H.R. 1
type FetchHR1Output struct {
	Body BillResponse
}

// RouteHandler holds dependencies for route handlers
type RouteHandler struct {
	billService *BillService
}

// NewRouteHandler creates a new RouteHandler with the given dependencies
func NewRouteHandler(billService *BillService) *RouteHandler {
	return &RouteHandler{billService: billService}
}

// --- Route Registration ---

// RegisterRoutes sets up all API routes with Huma (mock data fallback)
func RegisterRoutes(api huma.API) {
	// Health check
	huma.Get(api, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "healthy"
		resp.Body.Service = "deltagov-api"
		return resp, nil
	})

	// List all bills (mock data fallback)
	huma.Get(api, "/api/v1/bills", func(ctx context.Context, input *struct{}) (*ListBillsOutput, error) {
		bills := GetMockBills()
		resp := &ListBillsOutput{}
		resp.Body.Bills = mockBillsToBillResponses(bills)
		resp.Body.Total = len(bills)
		return resp, nil
	})
}

// RegisterRoutesWithService sets up all API routes with a real BillService
func RegisterRoutesWithService(api huma.API, handler *RouteHandler) {
	// Health check
	huma.Get(api, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "healthy"
		resp.Body.Service = "deltagov-api"
		return resp, nil
	})

	// Fetch H.R. 1 - The One Big Beautiful Bill
	huma.Register(api, huma.Operation{
		OperationID: "fetch-hr1",
		Method:      http.MethodPost,
		Path:        "/api/v1/bills/hr1/fetch",
		Summary:     "Fetch H.R. 1 (One Big Beautiful Bill)",
		Description: "Fetches H.R. 1 (119th Congress) from Congress.gov and stores all versions",
		Tags:        []string{"Bills"},
	}, func(ctx context.Context, input *struct{}) (*FetchHR1Output, error) {
		bill, err := handler.billService.FetchAndStoreHR1(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to fetch H.R. 1: " + err.Error())
		}
		return &FetchHR1Output{Body: *bill}, nil
	})

	// Get H.R. 1 directly (auto-fetch if not present)
	huma.Register(api, huma.Operation{
		OperationID: "get-hr1",
		Method:      http.MethodGet,
		Path:        "/api/v1/bills/hr1",
		Summary:     "Get H.R. 1 (One Big Beautiful Bill)",
		Description: "Returns H.R. 1 with all versions. Auto-fetches from Congress.gov if not cached.",
		Tags:        []string{"Bills"},
	}, func(ctx context.Context, input *struct{}) (*GetBillOutput, error) {
		bill, err := handler.billService.FetchAndStoreHR1(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get H.R. 1: " + err.Error())
		}
		return &GetBillOutput{Body: *bill}, nil
	})

	// List all bills
	huma.Register(api, huma.Operation{
		OperationID: "list-bills",
		Method:      http.MethodGet,
		Path:        "/api/v1/bills",
		Summary:     "List all bills",
		Description: "Returns all bills stored in the database",
		Tags:        []string{"Bills"},
	}, func(ctx context.Context, input *struct{}) (*ListBillsOutput, error) {
		bills, err := handler.billService.GetAllBills(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list bills: " + err.Error())
		}
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
		bill, err := handler.billService.GetBillByID(ctx, input.ID)
		if err != nil {
			return nil, huma.Error404NotFound("bill not found")
		}
		return &GetBillOutput{Body: *bill}, nil
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
		bill, err := handler.billService.GetBillWithVersions(ctx, input.ID)
		if err != nil {
			return nil, huma.Error404NotFound("bill not found")
		}
		resp := &GetBillVersionsOutput{}
		resp.Body.BillID = bill.ID
		resp.Body.Versions = bill.Versions
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
		diff, err := handler.billService.ComputeDiff(ctx, input.FromVersion, input.ToVersion)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to compute diff: " + err.Error())
		}
		return &ComputeDiffOutput{Body: *diff}, nil
	})
}

// mockBillsToBillResponses converts mock bills to BillResponse format
func mockBillsToBillResponses(mocks []MockBill) []BillResponse {
	responses := make([]BillResponse, len(mocks))
	for i, m := range mocks {
		id, _ := strconv.ParseUint(m.ID, 10, 32)
		responses[i] = BillResponse{
			ID:            uint(id),
			Title:         m.Title,
			Sponsor:       m.Sponsor,
			CurrentStatus: m.CurrentStatus,
		}
	}
	return responses
}
