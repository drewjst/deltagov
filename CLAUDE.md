# CLAUDE.md - DeltaGov Project Context

This file provides context for Claude Code when working on the DeltaGov project.

## Project Overview

DeltaGov is an open-source "Git for Government" platform designed to track, version, and visualize changes in U.S. legislative bills (starting with spending bills). It provides word-level "diffs" between bill versions and uses AI/ML to derive neutral, data-driven insights on policy shifts and spending changes.

**Core Value:** Transparency through version control. Treat laws like code.

**Business Model:** "Open Core" model. The engine and basic diffs are open-source (AGPLv3), while advanced predictive analytics and real-time lobbyist tracking will be proprietary "Premium" features.

**License:** GNU Affero General Public License v3.0 (AGPLv3)

## Technical Stack

| Layer          | Technology                  | Rationale                                              |
|----------------|-----------------------------|---------------------------------------------------------|
| Frontend       | Angular (v21)               | Enterprise-grade state management, strict typing        |
| Backend        | Go (Golang)                 | High-concurrency ingestion, fast string manipulation    |
| Database       | PostgreSQL (with JSONB)     | Hybrid relational/document storage for bill metadata    |
| Infrastructure | GCP (Cloud Run, Cloud SQL)  | Modular Monolith architecture                          |

## Monorepo Structure

```
/deltagov
├── /backend                 # Go module
│   ├── /cmd
│   │   ├── /api             # REST API (Fiber)
│   │   └── /ingestor        # Background worker (Congress.gov API poller)
│   ├── /internal            # Shared logic
│   │   ├── /api             # API route handlers
│   │   ├── /congress        # Congress.gov V3 API client
│   │   ├── /diff_engine     # Myers diff (go-udiff)
│   │   └── /models          # GORM structs (Bills, Versions, Deltas)
├── /frontend                # Angular workspace
│   └── /src/app             # Diff-viewer components & state services
├── /deployments             # Dockerfile & docker-compose.yml
└── LICENSE                  # GNU AGPLv3
```

## Key Domain Entities

### Bill
Metadata about a legislative bill:
- ID (Congress.gov bill identifier)
- Title
- Sponsor
- Current Status

### Version
A point-in-time snapshot of the bill text:
- Associated Bill ID
- ContentHash (SHA-256) to detect changes
- Raw text content
- Timestamp

### Delta
A structured JSON object representing the diff between two Versions:
- Insertions
- Deletions
- Unchanged text

## Core Workflows

### 1. Ingestion Pipeline
```
Poll Congress.gov API → Check if ContentHash is new → Store new Version in PostgreSQL
```

### 2. Diffing Process
```
Backend fetches two Versions → Computes text delta in Go (Myers diff) → Returns structured JSON
```

### 3. Visualization
```
Angular receives JSON delta → Renders side-by-side or inline diff view
```

## Development Commands

### Backend (Go)
```bash
cd backend

# Run API server
go run cmd/api/main.go

# Run ingestor worker
go run cmd/ingestor/main.go

# Run tests
go test ./...

# Build binaries
go build -o bin/api cmd/api/main.go
go build -o bin/ingestor cmd/ingestor/main.go
```

### Frontend (Angular)
```bash
cd frontend

# Install dependencies
npm install

# Development server
ng serve

# Build for production
ng build --configuration=production

# Run tests
ng test

# Run linting
ng lint
```

### Docker
```bash
cd deployments

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Coding Conventions

### Go Backend (Performance-First)
- **Allocations**: Pre-allocate slices/maps (`make([], 0, n)`) to avoid re-allocs.
- **Reuse**: Use `sync.Pool` for heavy objects like buffers in the `diff_engine`.
- **Concurrency**: Use `errgroup` for ingestion; always propagate `context.Context`.
- **Strings**: Use `strings.Builder` for bill content assembly.
- **Local Testing UI**: The API must serve an interactive **Scalar** or **Stoplight** UI at `/docs`.
- **Spec**: Always use **OpenAPI 3.1**.

### Angular Frontend (Spartan/Tailwind)
- **Library**: Standardize on **Spartan NG** (`@spartan-ng/ui-primitive`).
- **State**: Use **NGRX SignalStore** (`@ngrx/signals`) for feature state management; co-locate stores with components.
- **Styling**: Tailwind CSS exclusively; follow the "Helmet" (hlm) pattern for UI variants.
- **Architecture**: Standalone components only; use the `inject()` function.

### Git Workflow
- Main branch: `main`
- Feature branches: `feature/<description>`
- Bug fixes: `fix/<description>`
- Commit messages should be descriptive and reference issues when applicable

## External APIs

### Congress.gov API V3
- Documentation: https://api.congress.gov/
- Used for fetching bill metadata and text
- Requires API key (stored in environment variables)
- Rate limits apply

## Environment Variables

```bash
# Backend
CONGRESS_API_KEY=<your-api-key>
DATABASE_URL=postgres://user:pass@localhost:5432/deltagov
PORT=8080

# Frontend
API_BASE_URL=http://localhost:8080
```

## Architecture Notes

- **Modular Monolith:** Start simple, extract microservices later if needed
- **Content Hashing:** SHA-256 hash of bill text to efficiently detect changes
- **Myers Diff Algorithm:** Used for computing minimal edit distance between versions
- **JSONB Storage:** PostgreSQL JSONB for flexible delta storage while maintaining queryability
