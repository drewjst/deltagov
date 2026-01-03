# CLAUDE.md - DeltaGov Project Context

This file provides context for Claude Code when working on the DeltaGov project. It serves as a "constitution" for AI-assisted development, outlining project details, best practices, and guidelines to ensure consistency, performance, and maintainability. Keep this file concise, human-readable, and focused on high-impact guidance. Use it to enforce guardrails, document common patterns, and provide alternatives rather than just prohibitions.

For AI agents like Claude: Always think step-by-step before coding. Break down tasks into small iterations: plan, implement a minimal change, test, refine. Allocate extra "thinking time" for complex decisions—reason explicitly about trade-offs, performance implications, and alignment with best practices. If unsure, suggest alternatives or ask for clarification. Prioritize reuse, readability, and testing in every change.

## Project Overview

DeltaGov is an open-source "Git for Government" platform designed to track, version, and visualize changes in U.S. legislative bills (starting with spending bills). It provides word-level "diffs" between bill versions and uses AI/ML to derive neutral, data-driven insights on policy shifts and spending changes.

**Core Value:** Transparency through version control. Treat laws like code.

**Business Model:** "Open Core" model. The engine and basic diffs are open-source (AGPLv3), while advanced predictive analytics and real-time lobbyist tracking will be proprietary "Premium" features.

**License:** GNU Affero General Public License v3.0 (AGPLv3)

## Technical Stack

| Layer          | Technology                  | Rationale                                              |
|----------------|-----------------------------|---------------------------------------------------------|
| Frontend       | Angular (v21)               | Enterprise-grade state management, strict typing, zoneless change detection for performance |
| Backend        | Go (Golang)                 | High-concurrency ingestion, fast string manipulation, efficient for scalable APIs |
| Database       | PostgreSQL (with JSONB)     | Hybrid relational/document storage for bill metadata, efficient JSONB querying and indexing |
| Infrastructure | GCP (Cloud Run, Cloud SQL)  | Modular Monolith architecture, easy scaling and deployment |

## Monorepo Structure

```
/deltagov
├── /backend                 # Go module with layered architecture (controllers, services, models)
│   ├── /cmd
│   │   ├── /api             # REST API (Fiber) entry point
│   │   └── /ingestor        # Background worker (Congress.gov API poller)
│   ├── /internal            # Shared logic, separated by concerns
│   │   ├── /api             # API route handlers and middleware
│   │   ├── /congress        # Congress.gov V3 API client with dependency injection
│   │   ├── /diff_engine     # Myers diff (go-udiff), optimized for performance
│   │   ├── /db              # Database abstractions and utilities (e.g., GORM with connection pooling)
│   │   └── /models          # GORM structs (Bills, Versions, Deltas) with validation
├── /frontend                # Angular workspace with standalone components
│   └── /src/app             # Diff-viewer components, signal-based state services
├── /deployments             # Dockerfile, docker-compose.yml, Kubernetes manifests for scalability
├── /.ai                     # AI context files (e.g., subagents, prompts)
│   └── CLAUDE.local.md      # Local overrides (gitignore'd)
└── LICENSE                  # GNU AGPLv3
```

## Key Domain Entities

### Bill
Metadata about a legislative bill:
- ID (Congress.gov bill identifier)
- Title
- Sponsor
- Current Status
- (Best Practice: Use immutable structs; validate with embedded checks)

### Version
A point-in-time snapshot of the bill text:
- Associated Bill ID
- ContentHash (SHA-256) to detect changes efficiently
- Raw text content (stored as TEXT for large bills)
- Timestamp
- (Best Practice: Ensure immutability post-creation; use pure functions for hashing)

### Delta
A structured JSONB object representing the diff between two Versions:
- Insertions
- Deletions
- Unchanged text
- (Best Practice: Store in PostgreSQL JSONB with GIN indexes on common paths like 'insertions')

## Core Workflows

### 1. Ingestion Pipeline
```
Poll Congress.gov API → Check if ContentHash is new (using pure function) → Store new Version in PostgreSQL with transaction
```
- Use goroutines for concurrent polling; handle errors explicitly.

### 2. Diffing Process
```
Backend fetches two Versions → Computes text delta in Go (Myers diff) → Stores structured JSONB Delta
```
- Optimize with sync.Pool for buffers; pass dependencies explicitly.

### 3. Visualization
```
Angular receives JSON delta via API → Renders side-by-side or inline diff view using signals for reactivity
```
- Lazy-load components; use zoneless change detection for performance.

## Development Commands

### Backend (Go)
```bash
cd backend

# Run API server
go run cmd/api/main.go

# Run ingestor worker
go run cmd/ingestor/main.go

# Run tests with race detection
go test -race ./...

# Build binaries
go build -o bin/api cmd/api/main.go
go build -o bin/ingestor cmd/ingestor/main.go

# Lint and format
gofmt -w . && golangci-lint run
```

### Frontend (Angular)
```bash
cd frontend

# Install dependencies
npm install

# Development server with zoneless
ng serve --configuration=zoneless

# Build for production
ng build --configuration=production

# Run tests with Vitest
ng test

# Run linting (always before commits)
ng lint

# Typecheck
npm run typecheck
```

### Docker
```bash
cd deployments

# Start all services with scaling
docker-compose up -d --scale api=2

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Coding Conventions

### Go Backend (Performance-First, Functional Style)
- **Functional Programming**: Write in a functional style where possible—use pure functions, immutable data (e.g., pass copies of slices), explicit dependencies via interfaces. Avoid global state; inject loggers, DB clients, etc., as parameters. Use higher-order functions for abstractions like mapping over bills.
- **Allocations**: Pre-allocate slices/maps (`make([], 0, n)`) to avoid re-allocs.
- **Reuse**: Use `sync.Pool` for heavy objects like buffers in the `diff_engine`.
- **Concurrency**: Use `errgroup` for ingestion; always propagate `context.Context`. Leverage goroutines for scalability, but limit with channels to prevent overload.
- **Strings**: Use `strings.Builder` for bill content assembly.
- **Error Handling**: Always check errors; use multiple returns. Define custom errors with methods for clean reporting.
- **Architecture**: Follow layered structure (controllers → services → models). Use interfaces for decoupling; support mocking in tests.
- **Local Testing UI**: The API must serve an interactive **Scalar** or **Stoplight** UI at `/docs`.
- **Spec**: Always use **OpenAPI 3.1**. Validate DTOs with `validator`.
- **Best Practices**: Use connection pooling for DB; integrate caching (e.g., Redis) asynchronously. Profile with pprof; aim for stateless design for horizontal scaling.

### Angular Frontend (Spartan/Tailwind)
- **Components**: Use standalone components only; prefer `inject()` for dependencies. Break down into reusable parts for modularity.
- **State Management**: Use **NGRX SignalStore** (`@ngrx/signals`) for feature state; co-locate stores with components. Leverage signals for reactivity in forms (Signal Forms) and change detection.
- **Styling**: Tailwind CSS exclusively; follow the "Helmet" (hlm) pattern for UI variants. Use Angular Aria for accessible, headless components.
- **Type Checking**: Enable strict mode in `tsconfig.json`; avoid `any`; use type inference.
- **Performance**: Adopt zoneless change detection by default; lazy-load modules; use `defer` for heavy components.
- **Architecture**: Standalone APIs for simplicity; integrate AI tools (e.g., MCP Server) for code suggestions and migrations.
- **Best Practices**: Implement lazy loading; use Vitest for testing; run linting/typechecking iteratively during development.

### PostgreSQL Database
- **JSONB Usage**: Use JSONB for Deltas; normalize frequent fields (e.g., Bill ID) into columns. Hybrid modeling: relational for structured data, JSONB for flexible diffs.
- **Indexing**: Create GIN indexes on JSONB paths (e.g., `insertions` with `jsonb_path_ops` for containment). Use expression indexes for specific keys; partial indexes for subsets.
- **Querying**: Use containment operators (`@>`) for fast searches; cast extracted values (e.g., `(delta->>'count')::integer`). Avoid full document scans; use `jsonb_set` for targeted updates.
- **Validation**: Enforce with `CHECK` constraints (e.g., `IS JSON OBJECT`).
- **Performance**: Batch inserts; minimize updates to prevent bloat. Use migrations (e.g., Goose) for schema changes. Monitor with EXPLAIN; integrate full-text search for bill text.
- **Best Practices**: Prevent anti-patterns like duplicating data; use AWS enhancements (e.g., read replicas) on GCP equivalents.

## Git Workflow
- Main branch: `main`
- Feature branches: `feature/<description>`
- Bug fixes: `fix/<description>`
- Commit messages: Descriptive, reference issues (e.g., "feat: add diff visualization #123")
- Best Practice: Use TDD; commit small, iterative changes. Rebase before merging for clean history.

## Testing and CI/CD
- **Testing**: Aim for 80% coverage. Use TDD: write failing tests first, then code. Backend: `go test -race`; Frontend: `ng test` with Vitest.
- **Linting**: Always run before commits (Go: golangci-lint; Angular: ng lint).
- **CI/CD**: Use GitHub Actions or GCP pipelines. Lint, test, build on PRs; deploy to staging on merge. Include race detection and performance benchmarks.
- **Iteration**: Develop iteratively—prototype, test, refine. Allocate time for refactoring reusable components.

## Security Best Practices
- Store API keys in env vars or secrets manager (e.g., GCP Secrets).
- Use HTTPS; validate inputs to prevent injection.
- Rate-limit Congress.gov API calls; handle auth with middleware.
- Scan for vulnerabilities (e.g., `go vet`, npm audit).

## Performance Optimization
- **Backend**: Profile with pprof; use caching for diffs; goroutines for concurrency.
- **Frontend**: Zoneless CD; lazy loading; optimize signals.
- **Database**: Index heavily; batch operations; use connection pooling.
- Monitor: Integrate Prometheus/Grafana for metrics.

## External APIs

### Congress.gov API V3
- Documentation: https://api.congress.gov/
- Used for fetching bill metadata and text
- Requires API key (stored in environment variables)
- Rate limits apply—use backoff retries.

## Environment Variables

```bash
# Backend
CONGRESS_API_KEY=<your-api-key>
DATABASE_URL=postgres://user:pass@localhost:5432/deltagov
PORT=8080
REDIS_URL=redis://localhost:6379 # For caching

# Frontend
API_BASE_URL=http://localhost:8080
```

## Architecture Notes

- **Modular Monolith:** Start simple, extract microservices later if needed. Emphasize separation of concerns for reuse.
- **Content Hashing:** SHA-256 hash of bill text to efficiently detect changes.
- **Myers Diff Algorithm:** Used for computing minimal edit distance between versions; optimize for large texts.
- **JSONB Storage:** PostgreSQL JSONB for flexible delta storage while maintaining queryability; hybrid for performance.
- **Scalability:** Design for horizontal scaling; use stateless services.

## AI Assistance Guidelines
- **Approach Tasks**: Think step-by-step: Analyze requirements, plan architecture, code minimally, test, iterate. Consider reuse (e.g., extract functions) and context (e.g., performance in large bills).
- **Extra Thinking Time**: For complex changes, explicitly reason about alternatives, potential edge cases, and alignment with functional principles.
- **Iteration**: Build incrementally; suggest prototypes or PoCs. If code is suboptimal, propose refinements.
- **Common Pitfalls**: Avoid mutation in shared state; prefer explicit deps. For errors, use structured handling.
- **Subagents**: Use hierarchy (e.g., /backend/CLAUDE.md for Go-specific rules).
- **Workflow**: /init for setup; use checklists for tasks like "implement feature: 1. Plan, 2. Code, 3. Test".