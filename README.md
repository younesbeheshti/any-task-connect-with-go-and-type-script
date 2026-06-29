# TaskBridge (تسک‌بریج)

An escrow-based task marketplace connecting **Requesters** (task posters) with **Agents** (task
executors). The UI is in Persian (Farsi) with a right-to-left (RTL) layout.

The repository is a monorepo with two independent projects:

| Project | Stack |
|---------|-------|
| [`backend/`](backend/) | Go · Gin · GORM · PostgreSQL · Redis · RabbitMQ |
| [`front/`](front/) | TanStack Start · React 19 · TypeScript · Tailwind CSS v4 · shadcn/ui |

---

## Features

- **Auth & RBAC** — JWT access/refresh tokens with Redis-backed sessions; roles `ADMIN`,
  `REQUESTER`, `AGENT`. Agents register with city + national ID.
- **Tasks** — full lifecycle state machine with an escrow hold on creation.
- **Applications & assignment** — agents apply to open tasks; the requester reviews applicants and
  assigns exactly one. Duplicate applications and assigning an already-assigned task are rejected.
- **Wallet & payments** — mock instant top-up, escrow lock/release/refund, double-entry ledger.
- **Escrow payout handshake** — the requester verifies the work (escrow released to the agent), then
  the agent confirms receipt to reach the final `PAID` state.
- **Chat** — per-task conversations with attachments, read receipts, unread counts, and 5s polling.
- **Notifications** — in-app notifications (e.g. task cancelled, payment released, receipt confirmed).
- **File uploads** — local-disk storage with an authenticated download endpoint.
- **Admin, ratings, categories, cities, dashboard** — supporting bounded contexts.

### Task lifecycle

```
CREATED → OPEN → ASSIGNED → IN_PROGRESS → COMPLETED → WAITING_FOR_VERIFICATION → VERIFIED → PAID
           ↓         ↓            ↓              ↓
       CANCELLED  CANCELLED   CANCELLED      CANCELLED
```

- The requester verifies at `WAITING_FOR_VERIFICATION` → `VERIFIED` (escrow released to the agent).
- The assigned agent confirms receipt at `VERIFIED` → `PAID` (terminal).

---

## Getting started

### Prerequisites

- Go 1.22+
- Node.js 18+ (the frontend dev server runs on Node)
- Docker + Docker Compose (for Postgres, Redis, RabbitMQ)

### Backend

All commands run from `backend/`.

```bash
# Start infrastructure (Postgres, Redis, RabbitMQ) and run migrations
docker compose up postgres redis rabbitmq migrate

# Run the API server (http://localhost:8000)
go run ./cmd/api/

# Or run the full stack in Docker (builds the API image, starts everything)
docker compose up

# Tests
go test ./...
```

Migrations (require the `golang-migrate` CLI):

```bash
./scripts/migrate.sh up
./scripts/migrate.sh down
./scripts/migrate.sh create <name>
```

Swagger UI is served at `GET /swagger/*any` while the API runs.

### Frontend

All commands run from `front/`.

```bash
npm install      # install dependencies
npm run dev      # start the dev server (Vite)
npm run build    # production build
npm run lint     # ESLint
```

The frontend talks to the backend at `http://localhost:8000` by default (override with
`VITE_API_URL`).

---

## Architecture

The backend follows Clean Architecture with DDD. The dependency rule points inward; `domain/` has no
infrastructure imports.

| Layer | Path | Role |
|-------|------|------|
| Domain | `internal/*/domain/` | Entities, value objects, state machines |
| Repository | `internal/*/repository/` | Persistence interfaces (ports) |
| Service | `internal/*/service/` | Use-case interfaces + implementations |
| Infrastructure | `internal/*/infra/` | GORM repos, Redis/session adapters |
| Handler | `internal/*/handler/` | HTTP handlers + DTOs |
| Shared | `internal/common/` | RBAC, error types, events, cache keys |
| Packages | `pkg/` | Reusable infra (JWT, Redis, RabbitMQ, DB, logger, validator) |

Dependency wiring lives in `internal/bootstrap/app.go`; routes are registered in
`internal/api/router.go`.

**Bounded contexts:** `auth`, `user`, `task`, `application`, `wallet`, `payment`, `chat`,
`notification`, `rating`, `category`, `city`, `admin`, `file`.

The frontend uses TanStack Router (file-based routing under `src/routes/`), shadcn/ui components, and
Persian i18n utilities in `src/lib/`. The frontend–backend contract is documented in
[`front/docs/api-contracts.md`](front/docs/api-contracts.md).

### Conventions

- Base URL: `/v1` · Auth header: `Authorization: Bearer <jwt>`
- Currency: **Toman** as integer (no decimals)
- Dates: ISO-8601 in transport, Jalali calendar in the UI
- Error shape: `{ "error": { "code": string, "message": string, "fields"?: {} } }`
- Client-facing error messages are in **Persian**

---

## Configuration

Backend config loads in priority order: `configs/config.yaml` → `.env` → environment variables. Key
variables:

```
APP_PORT, APP_ENVIRONMENT
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE
REDIS_HOST, REDIS_PORT, REDIS_PASSWORD
RABBITMQ_URL, RABBITMQ_EXCHANGE
JWT_SECRET, JWT_ACCESS_TTL, JWT_REFRESH_TTL
STORAGE_LOCAL_DIR, STORAGE_MAX_SIZE
```
