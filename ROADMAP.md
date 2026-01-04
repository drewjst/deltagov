# DeltaGov Roadmap

This document tracks planned improvements, known issues, and architectural suggestions for the DeltaGov project.

## Current State

The codebase has solid foundations but is currently running on **mock data**. Real implementations exist for:
- Congress.gov API client (streaming JSON)
- Myers diff algorithm
- GORM database models
- Virtual scrolling UI

However, these components are not yet wired together.

---

## Critical: Data Flow Blockages

These issues prevent the application from functioning with real data.

### 1. Wire Database Connection

**Status:** Not Started
**Priority:** Critical
**Location:** `backend/cmd/api/main.go`

**Problem:** GORM models exist (`internal/models/bill.go`) but the database is never initialized.

**Action:**
```go
// Add to main.go
import "gorm.io/driver/postgres"

db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}
db.AutoMigrate(&models.Bill{}, &models.Version{}, &models.Delta{})
```

---

### 2. Replace Mock Data with Real Diff Engine

**Status:** Not Started
**Priority:** Critical
**Location:** `backend/internal/api/routes.go:127`

**Problem:** The `/diff/{from}/{to}` endpoint calls `GetMockDelta()` instead of `diff_engine.Compute()`.

**Action:**
- Inject database into route handlers
- Load version text content from DB
- Call `diff_engine.Compute(textA, textB, versionA, versionB)`
- Cache result in `deltas` table

---

### 3. Remove Frontend Mock Fallback

**Status:** Not Started
**Priority:** Critical
**Location:** `frontend/src/app/components/living-bill/living-bill.store.ts:88-111`

**Problem:** `diffLines` computed signal returns hardcoded mock data when no delta exists.

**Action:**
- Remove the `mockLines` fallback array
- Return empty array when `delta` is null
- Add loading state handling in the component

---

## Moderate: Architecture Improvements

These issues affect maintainability and code quality.

### 4. Eliminate Duplicate Type Definitions

**Status:** Not Started
**Priority:** Moderate
**Location:** `backend/internal/api/bills.go:22-38`

**Problem:** `FetchedBill` duplicates fields from `congress.Bill`. Multiple similar types exist:
- `congress.Bill` (API response)
- `api.FetchedBill` (duplicate)
- `api.MockBill` (mock data)
- `models.Bill` (database)

**Action:**
- Remove `FetchedBill` struct
- Create a single `BillDTO` for API responses
- Use `models.Bill` for database operations
- Convert between types at boundaries only

---

### 5. Add Service Layer

**Status:** Not Started
**Priority:** Moderate
**Location:** `backend/internal/api/routes.go`

**Problem:** Route handlers directly call mock functions with no abstraction layer.

**Action:**
Create `internal/service/bill_service.go`:
```go
type BillService interface {
    ListBills(ctx context.Context, offset, limit int) ([]models.Bill, int, error)
    GetBill(ctx context.Context, id string) (*models.Bill, error)
    GetVersions(ctx context.Context, billID string) ([]models.Version, error)
    ComputeDiff(ctx context.Context, billID, fromVersion, toVersion string) (*diff_engine.Delta, error)
}
```

---

### 6. Remove Unused Mutex

**Status:** Not Started
**Priority:** Low
**Location:** `backend/internal/congress/client.go:37`

**Problem:** `sync.RWMutex` is defined but never used.

**Action:**
- Either implement rate-limit tracking that uses the mutex
- Or remove it to reduce cognitive overhead

---

### 7. Implement Ingestor Worker

**Status:** Not Started
**Priority:** Moderate
**Location:** `backend/cmd/ingestor/main.go`

**Problem:** `pollBills()` function is a stub that does nothing.

**Action:**
```go
func pollBills(ctx context.Context, client *congress.Client, db *gorm.DB) error {
    result, err := client.FetchRecentBills(ctx, 50)
    if err != nil {
        return err
    }

    for _, bill := range result.Bills {
        // 1. Check if bill exists in DB
        // 2. Fetch bill text from Congress API
        // 3. Compute content hash
        // 4. If hash differs, create new Version
        // 5. Update Bill record
    }
    return nil
}
```

---

## Performance Optimizations

These issues may cause slowdowns at scale.

### 8. Add Diff Caching

**Status:** Not Started
**Priority:** High
**Location:** `backend/internal/api/routes.go:127`

**Problem:** Diff is recomputed on every request.

**Action:**
- Before computing, check `deltas` table for existing result
- After computing, store in `deltas` table
- Add cache invalidation when versions change

---

### 9. Lazy Load Bill Text

**Status:** Not Started
**Priority:** Medium
**Location:** `backend/internal/models/bill.go:49`

**Problem:** `TextContent` field loads entire bill text into memory.

**Action:**
- Consider streaming large texts
- Use external storage (S3/GCS) for bills over threshold
- Load text only when needed for diffing

---

### 10. Add Connection Pooling

**Status:** Not Started
**Priority:** Medium
**Location:** `backend/cmd/api/main.go`

**Problem:** No explicit connection pool configuration.

**Action:**
```go
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

---

### 11. Share HTTP Client

**Status:** Not Started
**Priority:** Low
**Location:** `backend/internal/congress/client.go:72`

**Problem:** Each `NewClient()` call creates a new `http.Client`.

**Action:**
- Use a shared HTTP client via singleton or dependency injection
- Enables HTTP/2 connection reuse

---

## Code Quality

Minor improvements for cleaner code.

### 12. Fix Double Tokenization in Word Diff

**Status:** Not Started
**Priority:** Low
**Location:** `backend/internal/diff_engine/diff.go:105`

**Problem:** Tokens are joined with `\n` then re-split by the diff algorithm.

**Action:**
- Diff token slices directly instead of joining/splitting

---

### 13. Add Pagination to List Bills

**Status:** Not Started
**Priority:** Medium
**Location:** `backend/internal/api/routes.go:76-82`

**Problem:** Returns all bills with no offset/limit.

**Action:**
```go
type ListBillsInput struct {
    Offset int `query:"offset" default:"0"`
    Limit  int `query:"limit" default:"20" maximum:"100"`
}
```

---

### 14. Move Mock Data to Separate Package

**Status:** Not Started
**Priority:** Low
**Location:** `frontend/src/app/components/living-bill/living-bill.store.ts:88-111`

**Problem:** Mock data is embedded in the signal store's computed function.

**Action:**
- Create `frontend/src/app/mocks/` directory
- Move mock data to dedicated mock service
- Keep store logic pure

---

## Frontend Tasks

### 15. Wire API Calls to Backend

**Status:** Not Started
**Priority:** Critical
**Location:** `frontend/src/app/components/living-bill/living-bill.ts`

**Problem:** Component initializes with mock data in `ngOnInit()` instead of fetching from API.

**Action:**
- Create `BillService` with HTTP client
- Call `GET /api/v1/bills/{id}/versions` on load
- Call `GET /api/v1/bills/{id}/diff/{from}/{to}` on version change

---

### 16. Implement Search Functionality

**Status:** Not Started
**Priority:** Medium
**Location:** `frontend/src/app/components/header/header.ts`

**Problem:** Search input exists but has no functionality.

**Action:**
- Add search endpoint to backend
- Debounce input in frontend
- Display search results dropdown

---

### 17. Add Error Handling UI

**Status:** Not Started
**Priority:** Medium
**Location:** `frontend/src/app/components/living-bill/`

**Problem:** No error states displayed to user.

**Action:**
- Use `store.error()` signal to show error messages
- Add retry button for failed requests
- Show loading skeletons during fetch

---

## Infrastructure

### 18. Add Database Migrations

**Status:** Not Started
**Priority:** High
**Location:** New directory `backend/migrations/`

**Problem:** Using `AutoMigrate` is not production-safe.

**Action:**
- Add Goose or golang-migrate
- Create initial migration from current models
- Document migration workflow

---

### 19. Set Up CI/CD

**Status:** Not Started
**Priority:** Medium
**Location:** `.github/workflows/`

**Action:**
- Add GitHub Actions workflow
- Run `go test -race ./...` on PRs
- Run `ng test` and `ng lint` on PRs
- Build Docker images on merge to main

---

### 20. Add Monitoring

**Status:** Not Started
**Priority:** Low
**Location:** `backend/cmd/api/main.go`

**Action:**
- Add Prometheus metrics endpoint
- Track request latency, error rates
- Add structured logging (zerolog/zap)

---

## Suggested Architecture

### Target State

```
cmd/api/main.go
    │
    ├── Initializes DB, Congress client, services
    │
    └── Injects into handlers
            │
            ▼
internal/service/
    ├── bill_service.go      ← Business logic
    └── diff_service.go      ← Caching layer
            │
            ▼
internal/repository/
    └── bill_repo.go         ← Database access
            │
            ▼
internal/models/             ← GORM structs (unchanged)
```

### Type Hierarchy

```
Congress.gov API
       │
       ▼
congress.Bill (raw API response)
       │
       ▼
models.Bill (database entity)
       │
       ▼
dto.BillResponse (API response to frontend)
```

---

## Completed

- [x] Project scaffolding
- [x] Fiber + Huma REST framework
- [x] Congress.gov API client with streaming JSON
- [x] Myers diff algorithm implementation
- [x] GORM model definitions
- [x] Angular 21 standalone components
- [x] NGRX SignalStore state management
- [x] Virtual scrolling for diff viewer
- [x] Spartan/Tailwind UI components
- [x] Docker infrastructure
- [x] OpenAPI/Scalar documentation

---

## Contributing

To work on any of these items:

1. Comment on the issue or create one referencing this roadmap item
2. Create a feature branch: `feature/<item-number>-short-description`
3. Submit a PR with tests
4. Reference this document in your PR description

---

*Last updated: January 2025*
