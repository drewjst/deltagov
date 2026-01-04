package congress

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	baseURL            = "https://api.congress.gov/v3"
	defaultTimeout     = 30 * time.Second
	defaultLimit       = 250 // Congress.gov max limit per request
	defaultPreallocCap = 250 // Pre-allocation capacity for bill slices
)

// Errors returned by the client.
var (
	ErrNoAPIKey      = errors.New("congress: API key is required")
	ErrInvalidStatus = errors.New("congress: unexpected status code")
	ErrRateLimited   = errors.New("congress: rate limit exceeded")
	ErrNotFound      = errors.New("congress: resource not found")
)

// Client is a thread-safe Congress.gov API V3 client.
// All methods are safe for concurrent use.
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string

	// mu protects any future mutable state (e.g., rate limit tracking)
	mu sync.RWMutex
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithAPIKey sets the Congress.gov API key.
func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithHTTPClient sets a custom HTTP client for the API requests.
// The provided client should be configured with appropriate timeouts.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithBaseURL overrides the default Congress.gov API base URL.
// Useful for testing with mock servers.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(url, "/")
	}
}

// New creates a new Congress.gov API client with the given API key.
// This is a convenience constructor for simple use cases.
func New(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: baseURL,
	}, nil
}

// NewClient creates a new Congress.gov API client with the given options.
// Returns an error if the API key is not provided.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: baseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	return c, nil
}

// Bill represents a legislative bill from Congress.gov API V3.
// Fields map to the /bill/{congress}/{billType} endpoint response.
// Note: Number is a string because some bill types use non-numeric identifiers.
type Bill struct {
	Congress                int           `json:"congress"`
	Type                    string        `json:"type"`
	Number                  string        `json:"number"`
	Title                   string        `json:"title"`
	OriginChamber           string        `json:"originChamber"`
	OriginChamberCode       string        `json:"originChamberCode"`
	UpdateDate              string        `json:"updateDate"`
	UpdateDateIncludingText string        `json:"updateDateIncludingText,omitempty"`
	URL                     string        `json:"url"`
	LatestAction            *LatestAction `json:"latestAction,omitempty"`
}

// LatestAction represents the most recent action on a bill.
type LatestAction struct {
	ActionDate string `json:"actionDate"`
	Text       string `json:"text"`
}

// BillsResponse represents the paginated API response for bills.
type BillsResponse struct {
	Bills      []Bill     `json:"bills"`
	Pagination Pagination `json:"pagination"`
}

// Pagination contains pagination metadata from the API response.
type Pagination struct {
	Count int    `json:"count"`
	Next  string `json:"next,omitempty"`
}

// FetchBillsResult contains the result of a FetchBills call.
type FetchBillsResult struct {
	Bills      []Bill
	TotalCount int
	HasMore    bool
}

// FetchBills retrieves bills for a specific congress and bill type.
// Uses streaming JSON decoding for memory efficiency.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - congress: The congress number (e.g., 118 for 118th Congress)
//   - billType: The bill type (e.g., "hr", "s", "hjres", "sjres")
//   - offset: Pagination offset (0-based)
//
// Returns FetchBillsResult with pre-allocated bill slice.
func (c *Client) FetchBills(ctx context.Context, congress int, billType string, offset int) (*FetchBillsResult, error) {
	url := fmt.Sprintf("%s/bill/%d/%s?api_key=%s&format=json&offset=%d&limit=%d",
		c.baseURL, congress, strings.ToLower(billType), c.apiKey, offset, defaultLimit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("congress: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("congress: failed to fetch bills: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	// Stream decode the response for memory efficiency
	// Pre-allocate the bills slice to avoid reallocations
	result := &FetchBillsResult{
		Bills: make([]Bill, 0, defaultPreallocCap),
	}

	decoder := json.NewDecoder(resp.Body)

	// Parse the opening brace
	if _, err := decoder.Token(); err != nil {
		return nil, fmt.Errorf("congress: failed to parse response start: %w", err)
	}

	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("congress: failed to parse key: %w", err)
		}

		switch key {
		case "bills":
			// Decode bills array using streaming
			if err := c.decodeBillsArray(decoder, result); err != nil {
				return nil, err
			}
		case "pagination":
			var pagination Pagination
			if err := decoder.Decode(&pagination); err != nil {
				return nil, fmt.Errorf("congress: failed to decode pagination: %w", err)
			}
			result.TotalCount = pagination.Count
			result.HasMore = pagination.Next != ""
		default:
			// Skip unknown fields
			var skip json.RawMessage
			if err := decoder.Decode(&skip); err != nil {
				return nil, fmt.Errorf("congress: failed to skip field %v: %w", key, err)
			}
		}
	}

	return result, nil
}

// decodeBillsArray streams the bills array from the JSON decoder.
// This avoids loading the entire array into memory at once.
func (c *Client) decodeBillsArray(decoder *json.Decoder, result *FetchBillsResult) error {
	// Consume the opening bracket of the array
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("congress: failed to parse bills array start: %w", err)
	}

	// Stream each bill object
	for decoder.More() {
		var bill Bill
		if err := decoder.Decode(&bill); err != nil {
			return fmt.Errorf("congress: failed to decode bill: %w", err)
		}
		result.Bills = append(result.Bills, bill)
	}

	// Consume the closing bracket of the array
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("congress: failed to parse bills array end: %w", err)
	}

	return nil
}

// checkResponse validates the HTTP response status code.
func (c *Client) checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusTooManyRequests:
		return ErrRateLimited
	default:
		return fmt.Errorf("%w: %d", ErrInvalidStatus, resp.StatusCode)
	}
}

// GetBillDetail fetches detailed information for a specific bill.
func (c *Client) GetBillDetail(ctx context.Context, congress int, billType string, billNumber int) (*Bill, error) {
	url := fmt.Sprintf("%s/bill/%d/%s/%d?api_key=%s&format=json",
		c.baseURL, congress, strings.ToLower(billType), billNumber, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("congress: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("congress: failed to fetch bill detail: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	// Response wraps bill in a "bill" key
	var wrapper struct {
		Bill Bill `json:"bill"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("congress: failed to decode bill detail: %w", err)
	}

	return &wrapper.Bill, nil
}

// GetBillText fetches the text versions available for a bill.
func (c *Client) GetBillText(ctx context.Context, congress int, billType string, billNumber int) ([]TextVersion, error) {
	url := fmt.Sprintf("%s/bill/%d/%s/%d/text?api_key=%s&format=json",
		c.baseURL, congress, strings.ToLower(billType), billNumber, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("congress: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("congress: failed to fetch bill text: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var wrapper struct {
		TextVersions []TextVersion `json:"textVersions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("congress: failed to decode text versions: %w", err)
	}

	return wrapper.TextVersions, nil
}

// TextVersion represents a text version of a bill.
type TextVersion struct {
	Date    string       `json:"date"`
	Type    string       `json:"type"`
	Formats []TextFormat `json:"formats"`
}

// TextFormat represents a specific format of bill text.
type TextFormat struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Appropriation/spending keywords for IsAppropriation check.
// Using lowercase for case-insensitive matching.
var appropriationKeywords = []string{
	"appropriation",
	"appropriations",
	"spending",
	"budget",
	"fiscal year",
	"continuing resolution",
	"omnibus",
}

// IsAppropriation checks if a bill title indicates it's an appropriations/spending bill.
// Uses optimized string matching for performance.
func IsAppropriation(title string) bool {
	if title == "" {
		return false
	}

	// Convert to lowercase once for all comparisons
	lower := strings.ToLower(title)

	// Use strings.Contains for each keyword - this is optimized in Go's stdlib
	// and uses efficient algorithms like Rabin-Karp for longer patterns
	for _, keyword := range appropriationKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}

	return false
}

// IsAppropriationFast is a faster variant that checks only the most common keywords.
// Use this in hot paths where performance is critical.
func IsAppropriationFast(title string) bool {
	if len(title) < 6 { // Shortest keyword is "budget"
		return false
	}

	lower := strings.ToLower(title)

	// Check most common patterns first (short-circuit evaluation)
	return strings.Contains(lower, "appropriation") ||
		strings.Contains(lower, "spending") ||
		strings.Contains(lower, "budget")
}

// FetchRecentBills retrieves the most recently updated bills from Congress.gov.
// This uses the /bill endpoint which returns bills sorted by updateDate descending.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - limit: Maximum number of bills to return (1-250)
//
// Returns FetchBillsResult with pre-allocated bill slice.
func (c *Client) FetchRecentBills(ctx context.Context, limit int) (*FetchBillsResult, error) {
	// Clamp limit to valid range
	if limit <= 0 {
		limit = 20
	}
	if limit > defaultLimit {
		limit = defaultLimit
	}

	url := fmt.Sprintf("%s/bill?api_key=%s&format=json&limit=%d&sort=updateDate+desc",
		c.baseURL, c.apiKey, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("congress: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("congress: failed to fetch recent bills: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	// Stream decode the response for memory efficiency
	result := &FetchBillsResult{
		Bills: make([]Bill, 0, limit),
	}

	decoder := json.NewDecoder(resp.Body)

	// Parse the opening brace
	if _, err := decoder.Token(); err != nil {
		return nil, fmt.Errorf("congress: failed to parse response start: %w", err)
	}

	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("congress: failed to parse key: %w", err)
		}

		switch key {
		case "bills":
			if err := c.decodeBillsArray(decoder, result); err != nil {
				return nil, err
			}
		case "pagination":
			var pagination Pagination
			if err := decoder.Decode(&pagination); err != nil {
				return nil, fmt.Errorf("congress: failed to decode pagination: %w", err)
			}
			result.TotalCount = pagination.Count
			result.HasMore = pagination.Next != ""
		default:
			var skip json.RawMessage
			if err := decoder.Decode(&skip); err != nil {
				return nil, fmt.Errorf("congress: failed to skip field %v: %w", key, err)
			}
		}
	}

	return result, nil
}
