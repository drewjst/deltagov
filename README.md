# DeltaGov

**Git for Government** — Track, version, and visualize changes in U.S. legislative bills.

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

## What is DeltaGov?

DeltaGov brings version control to legislation. Just as developers use Git to track code changes, DeltaGov tracks changes in legislative bills, providing:

- **Word-level diffs** between bill versions
- **Visual comparison** tools (side-by-side and inline views)
- **Change history** for any bill over time
- **AI-powered insights** on policy shifts and spending changes (coming soon)

Starting with U.S. spending bills, DeltaGov makes it easy to see exactly what changed, when, and by whom.

## Why DeltaGov?

Laws are the source code of society. Yet, while we track every semicolon in a mobile app, the trillions of dollars moving through spending bills are often hidden in 2,000-page PDFs. DeltaGov provides the Version Control System (VCS) that democracy is missing.

- **Journalists** investigating policy shifts
- **Researchers** studying legislative processes
- **Advocates** tracking amendments to bills they care about
- **Citizens** who want transparency in government

## Tech Stack

| Component      | Technology              |
|----------------|-------------------------|
| Frontend       | Angular 21              |
| UI Components  | Spartan NG, Tailwind CSS |
| Backend        | Go (Golang)             |
| Database       | PostgreSQL with JSONB   |
| Infrastructure | GCP Cloud Run       |

## Project Structure

```
/deltagov
├── /backend                        # Go API and ingestion workers
│   ├── /cmd
│   │   ├── /api                    # REST API entry point (Fiber + Huma)
│   │   └── /ingestor               # Background worker for Congress.gov polling
│   └── /internal
│       ├── /api                    # Route handlers and request/response types
│       ├── /config                 # Environment configuration loader
│       ├── /congress               # Congress.gov API V3 client (streaming JSON)
│       ├── /diff_engine            # Myers diff algorithm implementation
│       └── /models                 # GORM database models (Bill, Version, Delta)
├── /frontend                       # Angular 21 web application
│   ├── /src/app
│   │   ├── /components             # Standalone UI components
│   │   │   ├── /header             # App header with search
│   │   │   ├── /workspace          # Main layout container
│   │   │   ├── /living-bill        # Diff viewer with virtual scrolling
│   │   │   └── /insight-engine     # AI analysis pane
│   │   └── app.ts                  # Root component
│   └── /libs/ui                    # Spartan NG component library
├── /deployments
│   ├── Dockerfile                  # Multi-stage Go build
│   └── docker-compose.yml          # Local dev (PostgreSQL + Redis)
├── CLAUDE.md                       # AI development guidelines
├── ROADMAP.md                      # Suggested improvements and actions
└── LICENSE                         # GNU AGPLv3
```

## Architecture

### Component Overview

```mermaid
graph TB
    subgraph External["External Services"]
        CONGRESS["Congress.gov API V3"]
    end

    subgraph Backend["Backend (Go)"]
        API["REST API<br/>Fiber + Huma"]
        INGESTOR["Ingestor Worker"]
        CONGRESS_CLIENT["Congress Client"]
        DIFF_ENGINE["Diff Engine<br/>Myers Algorithm"]
        DB[(PostgreSQL)]
    end

    subgraph Frontend["Frontend (Angular 21)"]
        STORE["SignalStore"]
        LIVING["Living Bill<br/>Virtual Scroll"]
        INSIGHT["Insight Engine"]
    end

    CONGRESS --> CONGRESS_CLIENT
    INGESTOR --> CONGRESS_CLIENT
    CONGRESS_CLIENT --> DB
    API <--> DB
    API --> DIFF_ENGINE

    STORE <-->|HTTP| API
    STORE --> LIVING
    STORE --> INSIGHT
```

### Data Flow

```mermaid
flowchart LR
    subgraph Ingestion["Ingestion Pipeline"]
        A[Poll Congress.gov] --> B{New Content?}
        B -->|Yes| C[Compute Hash]
        C --> D[Store Version]
        B -->|No| E[Sleep]
    end

    subgraph Diffing["Diff Computation"]
        F[Fetch Versions] --> G[Myers Diff]
        G --> H[Store Delta]
    end

    subgraph Display["Visualization"]
        I[Load Delta] --> J[Virtual Scroll]
        J --> K[Color-coded Diff]
    end

    D --> F
    H --> I
```

### Request Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              REQUEST FLOW                                │
└─────────────────────────────────────────────────────────────────────────┘

  Browser                     API Server                    Database
  ═══════                     ══════════                    ════════

  GET /api/v1/bills ─────────────────>
                              Query bills table
                              <───────────────────────────────────
  <────────────────── JSON [{bills}]

  GET /api/v1/bills/{id}/versions ───>
                              Query versions by bill_id
                              <───────────────────────────────────
  <────────────────── JSON [{versions}]

  GET /api/v1/bills/{id}/diff/{a}/{b} >
                              Check delta cache
                              <───────────────────────────────────
                              If miss: compute Myers diff
                              Store delta for future
                              ────────────────────────────────────>
  <────────────────── JSON {segments, stats}
```

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 20.11+
- Angular CLI 21+
- PostgreSQL 15+
- Docker (optional)

### Quick Start with Docker

```bash
cd deployments
docker-compose up -d
```

### Manual Setup

**Backend:**
```bash
cd backend
go mod download
go run cmd/api/main.go
```

**Frontend:**
```bash
cd frontend
npm install
ng serve
```

The application will be available at `http://localhost:4200`.

### Environment Variables

```bash
# Required
CONGRESS_API_KEY=<your-congress-api-key>
DATABASE_URL=postgres://user:pass@localhost:5432/deltagov

# Optional
PORT=8080
```

Get a Congress.gov API key at: https://api.congress.gov/sign-up/

## How It Works

1. **Ingestion** — Polls Congress.gov API for bill updates
2. **Hashing** — Computes SHA-256 hash to detect text changes
3. **Diffing** — Uses Myers diff algorithm for minimal edit distance
4. **Visualization** — Renders changes in an intuitive diff viewer

## Ingestor Service

The ingestor is a background worker that fetches bills from Congress.gov and stores them in the database. It supports multiple ingestion modes and parallel processing.

### Ingestor CLI Flags

```bash
# Basic flags
--single-run              # Run once and exit (for Cloud Run Jobs)
--limit <n>               # Maximum bills to fetch (default: 50)

# Search-based ingestion
--search                  # Enable search mode (vs. recent bills mode)
--congress <n>            # Congress number to search (default: 119)
--type <type>             # Bill type: hr, s, hjres, sjres, hconres, sconres, hres, sres
--appropriations          # Only fetch appropriations/spending bills

# Performance
--parallel                # Use parallel processing for recent bills mode
--concurrency <n>         # Number of parallel workers (default: 5, max: 10)
```

### Usage Examples

```bash
# Fetch 50 most recent bills (default mode)
go run cmd/ingestor/main.go --single-run

# Fetch appropriations bills from the 119th Congress
go run cmd/ingestor/main.go --single-run --search --congress 119 --appropriations --limit 100

# Fetch House resolutions with parallel processing
go run cmd/ingestor/main.go --single-run --search --congress 119 --type hr --parallel --concurrency 8

# Continuous polling mode (for background service)
go run cmd/ingestor/main.go --search --appropriations
```

### Spending Bill Detection

The ingestor automatically identifies spending/appropriations bills using:
- **Policy Area filtering**: Uses Congress.gov's `policyArea` metadata
- **Keyword matching**: Detects titles containing "appropriation", "spending", "budget", "fiscal year", "continuing resolution", or "omnibus"

Bills marked as spending bills are flagged with `is_spending_bill=true` in the database for easy querying.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/bills` | List all tracked bills |
| GET | `/api/v1/bills/{id}` | Get bill details |
| GET | `/api/v1/bills/{id}/versions` | Get bill versions |
| GET | `/api/v1/bills/{id}/diff/{from}/{to}` | Compute diff between versions |
| GET | `/api/v1/lex` | Search bills with filters |
| GET | `/docs` | Interactive API documentation (Scalar) |
| GET | `/openapi.json` | OpenAPI 3.1 specification |

### Bill Search API (`/api/v1/lex`)

The Lex endpoint provides powerful search and filtering capabilities for legislative bills.

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `congress` | int | Filter by congress number (e.g., 118, 119). 0 = no filter |
| `sponsor` | string | Filter by sponsor name (case-insensitive partial match) |
| `query` | string | Search in bill title (case-insensitive partial match) |
| `type` | string | Filter by bill type: hr, s, hjres, sjres, hconres, sconres, hres, sres |
| `spending` | bool | Filter to only spending/appropriations bills |
| `limit` | int | Results per page (default: 20, max: 100) |
| `offset` | int | Pagination offset (default: 0) |

**Example Requests:**

```bash
# Search all bills (paginated)
curl "http://localhost:8080/api/v1/lex?limit=20"

# Search by title keyword
curl "http://localhost:8080/api/v1/lex?query=appropriation"

# Filter by congress and type
curl "http://localhost:8080/api/v1/lex?congress=119&type=hr"

# Get only spending bills
curl "http://localhost:8080/api/v1/lex?spending=true&limit=50"

# Combined filters with pagination
curl "http://localhost:8080/api/v1/lex?congress=119&sponsor=Johnson&limit=10&offset=20"
```

**Response Format:**

```json
{
  "bills": [
    {
      "id": 1,
      "congress": 119,
      "billNumber": 1,
      "billType": "hr",
      "title": "An act to provide for reconciliation...",
      "sponsor": "Rep. Smith",
      "originChamber": "House",
      "currentStatus": "Became Public Law",
      "updateDate": "2025-12-30T17:32:50Z"
    }
  ],
  "total": 150,
  "limit": 20,
  "offset": 0
}

## Contributing

We welcome contributions! Please see our contributing guidelines (coming soon).

### Development

For detailed development instructions and coding conventions, see [CLAUDE.md](./CLAUDE.md).

For planned improvements and known issues, see [ROADMAP.md](./ROADMAP.md).

## License

DeltaGov is open source under the [GNU Affero General Public License v3.0](LICENSE).

**Open Core Model:** The core diffing engine is open source. Premium features (predictive analytics, real-time lobbyist tracking) will be proprietary.

## Roadmap

- [x] Project scaffolding
- [x] Congress.gov API client (streaming JSON)
- [x] Diff engine (Myers algorithm)
- [x] Frontend layout with virtual scrolling
- [x] Database integration (PostgreSQL with GORM)
- [x] Ingestor worker with search-based ingestion
- [x] Batch processing with errgroup concurrency
- [x] GORM upsert for bill metadata updates
- [x] Appropriations/spending bill detection
- [x] Bill search API (`/api/v1/lex`) with filtering and pagination
- [x] Lex page (Angular) with search UI and bill listing
- [ ] Wire real data flow (replace mocks)
- [ ] AI-powered change summaries
- [ ] Premium analytics features

See [ROADMAP.md](./ROADMAP.md) for detailed improvement plans.

## Contact

Questions or feedback? Open an issue on GitHub.

---

*Transparency through version control. Treat laws like code.*
