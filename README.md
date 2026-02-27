# FSWithGo

A full-stack issue tracker built with a **Go REST API** backend and a React/Next.js frontend. The backend is designed around clean architecture principles with a focus on type safety, security, and maintainability — using only the Go standard library for routing and middleware.

## Backend Highlights

### Zero-Dependency HTTP Layer

The entire HTTP layer is built on Go 1.22+'s enhanced `net/http` — no third-party routers or frameworks. Routes use the new method-pattern syntax (`"GET /v1/issues/{id}"`) and middleware is composed through a custom `CreateStack()` function.

### Middleware Pipeline

Six middleware layers execute on every request in order:

| Middleware | Purpose |
|---|---|
| **CORS** | Configurable origin, handles OPTIONS preflight with 204 |
| **Recoverer** | Catches panics, logs full stack trace, returns 500 |
| **RequestID** | Extracts or generates UUID v4, echoes via `X-Request-ID` header |
| **RealIP** | Resolves client IP through trusted proxies (Cloudflare, Nginx, X-Forwarded-For) with private IP rejection |
| **Logging** | Logs `<status> <method> <path> <duration>` for every request |
| **Timeout** | 25-second deadline with buffered response — returns 504 on expiry |

### Authentication & Authorization

- **JWT (HS256)** with 7-day expiry, generated at login and registration
- **Bitfield permissions** baked into the JWT payload:

  | Permission | Bit | Value |
  |---|---|---|
  | Read | 0 | 1 |
  | Write | 1 | 2 |
  | Admin | 2 | 4 |

- **Three-tier auth middleware**: `GlobalAuth` (extracts token globally) → `RequiredAuth` (guards protected routes) → `RequiredAdmin` (guards admin routes)
- Passwords hashed with **bcrypt** (cost 10), never serialized in responses (`json:"-"`)

### API Endpoints

**Public**

| Method | Path | Description |
|---|---|---|
| `GET` | `/v1/health` | Liveness probe |
| `POST` | `/v1/users/register` | Create account, returns user + JWT |
| `POST` | `/v1/users/login` | Authenticate, returns user + JWT |
| `GET` | `/v1/issues/{id}` | Fetch single issue (with author name via JOIN) |

**Protected** (valid JWT required)

| Method | Path | Description |
|---|---|---|
| `GET` | `/v1/issues` | List authenticated user's issues |
| `POST` | `/v1/issues` | Create a new issue |
| `DELETE` | `/v1/issues/{id}` | Delete an issue |
| `PATCH` | `/v1/issues/{id}/status` | Update issue status |

**Admin** (JWT + admin permission bit)

| Method | Path | Description |
|---|---|---|
| `GET` | `/v1/admin/stats` | Total users, total issues, per-status breakdown |

### Repository Pattern & Anti-Corruption Layer

The data layer follows a clean separation of concerns:

```
Handler → Domain Interface → Store Implementation → sqlc Queries → PostgreSQL
```

- **`Storage` struct** holds interface fields (`Issues`, `Users`, `Admin`) — fully mockable for testing without any third-party library
- **Domain types** (`store.Issue`, `store.User`) are completely separate from sqlc-generated types, with explicit translation between `pgtype.Timestamptz` ↔ `time.Time`
- **Sentinel errors** (`store.ErrNotFound`) provide consistent error handling across the store layer

### Type-Safe SQL with sqlc

SQL queries are written as annotated `.sql` files and compiled into type-safe Go code at build time. The generated code uses **pgx/v5 native mode** (binary protocol, not `database/sql`) for direct PostgreSQL type support including UUIDs, enums, and timestamps.

Type overrides map PostgreSQL `uuid` → `google/uuid.UUID` and the `status_type` enum → Go `string`.

### Request Validation

- Request bodies capped at **1 MB**
- Unknown JSON fields are rejected (`DisallowUnknownFields`)
- Descriptive error messages for malformed JSON (syntax errors, wrong types, empty body, oversized payload)
- Field-level validation with structured `422` responses (`{"error": {"field": "message"}}`)
- Email validated against RFC-5321 regex; passwords require 8–20 chars with uppercase, lowercase, digit, and special character
- Duplicate email registration returns `409 Conflict` (catches PostgreSQL unique violation code `23505`)

### Database

**PostgreSQL 17** with pgx/v5 connection pooling:
- Configurable max connections (`DB_MAX_CONNS`, default 25) and idle timeout (`DB_MAX_IDLE_MINS`, default 15 min)
- Startup health check with 5-second timeout
- Schema managed through **golang-migrate** with numbered up/down SQL pairs

**Schema:**

```
users                          issues
─────                          ──────
id         UUID PK             id          BIGSERIAL PK
email      TEXT UNIQUE         user_id     UUID FK → users.id
password_hash TEXT             title       TEXT
name       TEXT                description TEXT
permissions INTEGER (default 3) status    status_type ENUM
created_at TIMESTAMPTZ         created_at  TIMESTAMPTZ
updated_at TIMESTAMPTZ         updated_at  TIMESTAMPTZ
```

### Docker

Multi-stage Dockerfile produces a minimal Alpine image:

1. **Builder stage** (`golang:1.25-alpine`): compiles API server + migration runner
2. **Runtime stage** (`alpine:3.21`): runs as non-root `appuser` (UID 1000)

The entrypoint script automatically converts `postgres://` → `pgx5://` and runs pending migrations before starting the server.

Docker Compose orchestrates three services: PostgreSQL (with healthcheck), the Go API, and the frontend dev server.

## Quick Start

```bash
# Start everything with Docker
docker compose up

# Or run locally:
# 1. Start Postgres (Docker)
docker compose up db

# 2. Run migrations
DATABASE_URL="pgx5://postgres:postgres@localhost:5433/fswithgo?sslmode=disable" go run ./cmd/migrate -direction up

# 3. Start the API
DATABASE_URL="postgres://postgres:postgres@localhost:5433/fswithgo?sslmode=disable" \
JWT_SECRET="your-secret" \
CORS_ORIGIN="http://localhost:3000" \
go build -o bin/main ./cmd/api && ./bin/main

# 4. Start the frontend
cd web && pnpm install && pnpm dev
```

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `ADDR` | `:3000` | HTTP listen address |
| `DATABASE_URL` | `postgres://...localhost:5432/fswithgo` | PostgreSQL connection string |
| `DB_MAX_CONNS` | `25` | Connection pool max size |
| `DB_MAX_IDLE_MINS` | `15` | Connection idle timeout (minutes) |
| `JWT_SECRET` | — | HMAC signing key for JWTs |
| `CORS_ORIGIN` | — | Allowed CORS origin |

## Tech Stack

**Backend:** Go · net/http · pgx/v5 · sqlc · golang-migrate · bcrypt · JWT (HS256)

**Frontend:** Next.js · React 19 · TypeScript · TanStack Query · Zod · Tailwind CSS · shadcn/ui

**Infrastructure:** PostgreSQL 17 · Docker · Nginx (production frontend)
