# Kopelko Dating App Backend

Production-grade Dating App REST API backend written in Go, architected as a **Modular Monolith** using **Clean Architecture** (*handler $\to$ usecase $\to$ repository*), built with **Echo v4** and **jackc/pgx/v5** connection pool.

---

## 🌟 Key Features & Engineering Highlights

- **Modular Monolith Architecture**: Strict domain boundary separation (`identity`, `profile`, `swipe`, `subscription`) inside `internal/modules/` with interface-based inter-module communication.
- **High-Performance PostgreSQL Pool**: Direct raw SQL queries with `pgxpool.Pool` (no heavy ORM overhead, full control over query performance and execution plans).
- **Transaction Atomicity & Concurrency Safety**:
  - Atomic user registration + initial profile creation.
  - Concurrency-safe swiping, atomic quota checking, and mutual matchmaker inside PostgreSQL transactions (`pgx.Tx`).
- **Discovery Feed Engine**: Random candidate feed filtered by gender and age range preferences while excluding previously viewed profiles and enforcing daily view limits.
- **Security Hardening**:
  - Anti brute-force rate limiter on authentication endpoints (`POST /api/register` & `POST /api/login`) based on client IP.
  - Standard web security headers (`nosniff`, `DENY`, `XSS Protection`, `HSTS`).
  - Strict input payload validation (`go-playground/validator`).
  - Passwords securely hashed with bcrypt.
- **Observability & Health Probes**:
  - Structured logging with Go's `log/slog` (JSON format in production, colored text in development).
  - Traceable `X-Request-ID` middleware injected to every request context.
  - Kubernetes / Coolify-ready health endpoints: `/health` (Liveness) & `/health/ready` (Readiness DB ping).
- **Interactive Documentation**: Swagger / OpenAPI 2.0 UI embedded at `/swagger/index.html`.
- **CI/CD & Code Quality**: Fully automated GitHub Actions workflow with `golangci-lint`, race condition detection (`-race`), and build checks.

---

## 🏗️ Project Structure

```
.
├── cmd
│   └── api
│       └── main.go           # Application entrypoint & dependency wiring
├── internal
│   ├── modules
│   │   ├── identity          # User Registration, Authentication & JWT
│   │   │   ├── delivery/http
│   │   │   ├── domain
│   │   │   ├── repository
│   │   │   └── usecase
│   │   ├── profile           # User Profiles & Discovery Feed
│   │   ├── subscription      # Premium Features & Quota Unlocking
│   │   └── swipe             # Swiping Logic & Mutual Matchmaking
│   └── platform              # Shared Core Platform
│       ├── config            # Environment loader
│       ├── database          # pgxpool connection & Transactor
│       ├── http              # Unified API response formatter
│       ├── logger            # Structured slog with context support
│       ├── middleware        # JWT Auth, Rate Limiter, Security & Request-ID
│       ├── migrator          # Programmatic schema & seeder runner
│       └── token             # JWT generation & verification
├── databases
│   ├── migrations            # PostgreSQL DDL schema & composite indexes
│   └── seeders               # Idempotent dummy data seeder
├── docs                      # Auto-generated Swagger specifications
├── Dockerfile                # Lightweight multi-stage build container
├── docker-compose.yml        # PostgreSQL 16 + API service stack
└── Makefile                  # Development task automation
```

---

## 🚀 Quickstart Guide

### Prerequisites
- Go 1.25+
- Docker & Docker Compose (optional, for containerized run)
- PostgreSQL 14+ (if running bare-metal)

### 1. Run with Docker Compose (Fastest)

```bash
docker compose up -d --build
```
The database will be automatically initialized with schemas and seeders.
- API is accessible at: `http://localhost:8080`
- Swagger UI documentation: `http://localhost:8080/swagger/index.html`
- Healthcheck endpoint: `http://localhost:8080/health/ready`

### 2. Run Locally

1. Setup environment file:
   ```bash
   cp .env.example .env # or configure environment variables
   ```
2. Run database migration and initial seed data:
   ```bash
   make migrate
   make seed
   ```
3. Start server:
   ```bash
   make run
   ```

---

## 🧪 Testing & Code Quality

Run tests with Go's race detector enabled:
```bash
make test
```

Generate HTML coverage report:
```bash
make test-cover
# Report saved to coverage.html
```

Run static analysis & linter:
```bash
make lint
```

---

## 📖 API Endpoints Reference

| Method | Endpoint | Auth | Description |
|---|---|:---:|---|
| `GET` | `/health` | No | Liveness probe |
| `GET` | `/health/ready` | No | Readiness probe (Verifies PostgreSQL connection) |
| `GET` | `/swagger/*` | No | Swagger OpenAPI 2.0 UI |
| `POST` | `/api/register` | Rate-limited | Register new account and profile |
| `POST` | `/api/login` | Rate-limited | Login and obtain JWT Bearer token |
| `GET` | `/api/users/profiles/me` | Bearer | Get authenticated user profile with active badges |
| `PUT` | `/api/users/profiles/me` | Bearer | Update profile bio, photos, interests, etc. |
| `GET` | `/api/users/profiles/random` | Bearer | Discovery feed (supports `?gender=&min_age=&max_age=`) |
| `POST` | `/api/users/swipes/:target_user_id` | Bearer | Swipe like / pass candidate |
| `GET` | `/api/users/matches` | Bearer | List mutual matches with profile details |
| `POST` | `/api/users/subscriptions` | Bearer | Subscribe to premium feature (`no_swipe_quota`, `verified_label`) |

---

## 📜 License

MIT License.
