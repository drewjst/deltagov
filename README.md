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
│   ├── cmd/
│   │   ├── api/                    # Web API entry point
│   │   └── ingestor/               # Background worker entry point
│   └── internal/
│       ├── api_client/             # Congress.gov API wrapper
│       ├── database/               # PostgreSQL/GORM layer
│       ├── models/                 # Shared Go structs
│       └── processor/              # Diff engine logic
├── /frontend                       # Angular web application
│   └── src/app/
│       ├── components/             # UI components (diff viewer, bill list, etc.)
│       ├── services/               # API communication & business logic
│       └── models/                 # TypeScript interfaces
├── /deployments
│   ├── /docker
│   │   ├── api.Dockerfile          # Multi-stage Distroless build
│   │   ├── ingestor.Dockerfile     # Same binary, different entrypoint
│   │   └── frontend.Dockerfile     # Nginx-based static host
│   ├── docker-compose.yml          # For local dev orchestration
│   └── cloud-run.sh                # Helper script for gcloud CLI deploys
└── LICENSE                         # GNU AGPLv3
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

## Contributing

We welcome contributions! Please see our contributing guidelines (coming soon).

### Development

For detailed development instructions and coding conventions, see [CLAUDE.md](./CLAUDE.md).

graph TD
    subgraph "External Data"
        C[Congress.gov API V3]
    end

    subgraph "DeltaGov Backend (Go)"
        I[Ingestor Worker] -->|Poll & Hash| H{Content Change?}
        H -->|Yes| DB[(PostgreSQL)]
        H -->|No| S[Sleep/Wait]
        D[Diff Engine] <-->|Myers Algorithm| DB
        A[REST API] <--> DB
    end

    subgraph "Frontend (Angular 21)"
        F[SignalStore] <--> A
        F --> V[Spartan UI / Diff Viewer]
    end

## License

DeltaGov is open source under the [GNU Affero General Public License v3.0](LICENSE).

**Open Core Model:** The core diffing engine is open source. Premium features (predictive analytics, real-time lobbyist tracking) will be proprietary.

## Roadmap

- [x] Project scaffolding
- [ ] Congress.gov API integration
- [ ] Basic diff engine (Myers algorithm)
- [ ] Bill version storage
- [ ] Diff visualization UI
- [ ] AI-powered change summaries
- [ ] Premium analytics features

## Contact

Questions or feedback? Open an issue on GitHub.

---

*Transparency through version control. Treat laws like code.*
