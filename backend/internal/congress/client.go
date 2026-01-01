package congress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL = "https://api.congress.gov/v3"
)

// Client wraps the Congress.gov API
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// Bill represents a legislative bill from Congress.gov
type Bill struct {
	Congress       int    `json:"congress"`
	Type           string `json:"type"`
	Number         int    `json:"number"`
	Title          string `json:"title"`
	OriginChamber  string `json:"originChamber"`
	UpdateDate     string `json:"updateDate"`
	LatestActionDate string `json:"latestAction,omitempty"`
	URL            string `json:"url"`
}

// BillsResponse represents the API response for bills
type BillsResponse struct {
	Bills      []Bill `json:"bills"`
	Pagination struct {
		Count int    `json:"count"`
		Next  string `json:"next,omitempty"`
	} `json:"pagination"`
}

// NewClient creates a new Congress.gov API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetRecentBills fetches recently updated bills
func (c *Client) GetRecentBills(ctx context.Context) ([]Bill, error) {
	url := fmt.Sprintf("%s/bill?api_key=%s&format=json&limit=20&sort=updateDate+desc", baseURL, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bills: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result BillsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Bills, nil
}

// GetBillText fetches the full text of a specific bill version
func (c *Client) GetBillText(ctx context.Context, congress int, billType string, billNumber int) (string, error) {
	url := fmt.Sprintf("%s/bill/%d/%s/%d/text?api_key=%s&format=json",
		baseURL, congress, billType, billNumber, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch bill text: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// TODO: Parse text versions response and fetch actual text content
	return "", nil
}
