# TaskBridge Backend — Architecture

## Overview

TaskBridge is an escrow-based marketplace connecting **Requesters** (people who need tasks done) with **Agents** (people who perform tasks physically). The backend is a Go monolith organized by bounded contexts, following Clean Architecture, DDD, and event-driven patterns.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Layer                                    │
│         Web App (React)  │  Mobile (future)  │  Admin Dashboard             │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │ HTTPS / WSS
┌───────────────────────────────────▼─────────────────────────────────────────┐
│                           Delivery Layer (Gin)                               │
│  HTTP Handlers  │  WebSocket Hub  │  Middleware (JWT, RBAC, Rate Limit)     │
│  /api/v1/*      │  /ws            │  Request Validation  │  Swagger          │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
┌───────────────────────────────────▼─────────────────────────────────────────┐
│                         Application Layer (Services)                         │
│  AuthService  TaskService  ApplicationService  WalletService  PaymentService │
│  ChatService  NotificationService  ReviewService  AdminService             │
└───────┬─────────────────┬──────────────────┬────────────────┬───────────────┘
        │                 │                  │                │
        ▼                 ▼                  ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐
│  Repository  │  │    Redis     │  │   RabbitMQ   │  │  External Gateways   │
│  (GORM/PG)   │  │  Cache/RL/   │  │  Events/Jobs │  │  SMS, Email, Payment │
│              │  │  Sessions    │  │              │  │  Object Storage      │
└──────────────┘  └──────────────┘  └──────────────┘  └──────────────────────┘
        │
        ▼
┌──────────────┐
│  PostgreSQL  │
└──────────────┘
```

## Layer Responsibilities

| Layer | Location | Responsibility |
|-------|----------|----------------|
| **Domain** | `internal/*/domain/` | Entities, value objects, state machines, domain errors |
| **Repository** | `internal/*/repository/` | Persistence interfaces (ports) |
| **Service** | `internal/*/service/` | Use-case interfaces and orchestration contracts |
| **Infrastructure** | `internal/*/infra/` | GORM repos, Redis, RabbitMQ adapters (Phase 2+) |
| **Delivery** | `internal/*/handler/` | HTTP handlers, DTOs, route registration (Phase 2+) |
| **Shared** | `internal/common/`, `pkg/` | Cross-cutting utilities |

Dependency rule: **outer layers depend on inner layers**. Domain has zero infrastructure imports.

## Bounded Contexts

| Context | Package | Core Aggregates |
|---------|---------|-----------------|
| Auth | `internal/auth` | Session, Token, OTP |
| User | `internal/user` | User, Verification |
| Task | `internal/task` | Task, TaskTimeline |
| Application | `internal/application` | Application |
| Wallet | `internal/wallet` | Wallet |
| Payment | `internal/payment` | Transaction, Escrow |
| Chat | `internal/chat` | Chat, Message |
| Notification | `internal/notification` | Notification |
| Rating | `internal/rating` | Review |
| Category | `internal/category` | Category |
| City | `internal/city` | City |
| Admin | `internal/admin` | Metrics, Moderation |

## API Versioning & Response Envelope

- Base path: `/api/v1`
- Success: `{ "success": true, "message": "...", "data": {} }`
- Error: `{ "success": false, "message": "...", "errors": {} }`

Frontend contracts (`front/docs/api-contracts.md`) use camelCase JSON and a compatible error shape. Handlers will map domain errors to the envelope; optional `code` field can be added in `errors` for client routing.

## Authentication & Authorization

### JWT Strategy

| Token | TTL | Storage | Purpose |
|-------|-----|---------|---------|
| Access | 15 min | Client | API authorization |
| Refresh | 7 days | HttpOnly cookie + Redis session | Token rotation |

Claims: `sub` (user UUID), `role`, `jti`, `exp`, `iat`.

Logout blacklists `jti` in Redis until access token expiry.

### RBAC

```
Role        Permissions
────────────────────────────────────────────────────────────
ADMIN       Full system access, moderation, analytics
REQUESTER   CRUD own tasks, accept apps, verify, wallet, chat
AGENT       Browse tasks, apply, execute assigned tasks, wallet, chat
```

Permission checks use middleware composition:

```
JWTAuth → RoleGuard(roles...) → PermissionGuard(perm...) → Handler
```

## Task State Machine

Internal domain states (UPPER_SNAKE) map to API values (snake_case):

| Domain State | API Status | Description |
|--------------|------------|-------------|
| `CREATED` | `posted` | Task record created, not yet published |
| `OPEN` | `awaiting_applicants` | Published, accepting applications |
| `ASSIGNED` | `accepted` | Agent selected |
| `IN_PROGRESS` | `in_progress` | Agent started work |
| `COMPLETED` | `completed` | Agent marked complete |
| `WAITING_FOR_VERIFICATION` | `awaiting_verification` | Awaiting requester approval |
| `VERIFIED` | *(internal)* | Requester verified; triggers payment |
| `PAID` | `paid` | Escrow released to agent |
| `CANCELLED` | `cancelled` | Cancelled with refund rules |

Valid transitions are enforced in `Task.CanTransitionTo()` — invalid transitions return `ErrInvalidTaskTransition`.

## Escrow Payment Flow

```
┌─────────────┐     Create Task      ┌──────────────────────────────────┐
│  Requester  │ ──────────────────►  │  BEGIN TRANSACTION               │
│   Wallet    │                      │  available -= budget + fee       │
└─────────────┘                      │  locked += budget + fee          │
                                     │  INSERT task, INSERT tx (LOCK)   │
                                     │  COMMIT                          │
                                     └──────────────────────────────────┘

┌─────────────┐     Verify Task      ┌──────────────────────────────────┐
│  Requester  │ ──────────────────►  │  BEGIN TRANSACTION               │
│   Wallet    │                      │  requester.locked -= held          │
│  Agent      │                      │  agent.available += budget       │
│   Wallet    │                      │  platform fee ledger (optional)  │
└─────────────┘                      │  task → PAID, INSERT tx (RELEASE)│
                                     │  COMMIT → publish PaymentReleased│
                                     └──────────────────────────────────┘
```

All wallet mutations use `SELECT ... FOR UPDATE` on wallet rows within a single PostgreSQL transaction.

## Redis Usage

| Key Pattern | TTL | Purpose |
|-------------|-----|---------|
| `jwt:blacklist:{jti}` | access TTL | Revoked tokens |
| `session:{userID}:{sessionID}` | 7d | Refresh session metadata |
| `ratelimit:{ip}:{route}` | 1m | Rate limiting (sliding window) |
| `online:{userID}` | 5m | WebSocket presence |
| `task:{id}` | 5m | Task detail cache (invalidate on write) |
| `notifications:{userID}` | 2m | Unread notification count cache |

## RabbitMQ Events

Exchange: `taskbridge.events` (topic)

| Routing Key | Publisher | Consumers |
|-------------|-----------|-----------|
| `task.created` | TaskService | Notification, Analytics |
| `task.assigned` | ApplicationService | Notification, WebSocket |
| `task.started` | TaskService | Notification, WebSocket |
| `task.completed` | TaskService | Notification, WebSocket |
| `task.verified` | TaskService | PaymentService |
| `payment.released` | PaymentService | Notification, WebSocket |
| `application.submitted` | ApplicationService | Notification, WebSocket |
| `notification.created` | NotificationService | WebSocket |
| `review.created` | RatingService | UserService (rating aggregate) |

## WebSocket Realtime

Hub pattern at `/ws`. Clients authenticate via JWT query param or first message.

Events pushed to connected clients:

- `NewApplication`
- `TaskAssigned`
- `TaskCompleted`
- `PaymentReleased`
- `NewMessage`
- `NotificationCreated`

## Observability

- **Logging**: Zap structured JSON (`request_id`, `user_id`, `latency_ms`)
- **Health**: `GET /health` (liveness), `GET /ready` (DB + Redis + RabbitMQ)
- **Metrics**: `GET /metrics` (Prometheus-compatible, Phase 2)

## Security Checklist

- bcrypt password hashing (cost 12)
- Parameterized queries via GORM
- Input validation via go-playground/validator
- Rate limiting on auth and write endpoints
- Secure headers (HSTS, X-Content-Type-Options, etc.)
- CORS restricted to configured origins

## Deployment Topology (Docker Compose)

```
api ──┬── postgres
      ├── redis
      ├── rabbitmq
      └── worker (event consumers)
```

## Phase Roadmap

| Phase | Scope |
|-------|-------|
| **1 (current)** | Architecture, schema, models, interfaces |
| **2** | Infrastructure adapters, migrations, DI wiring |
| **3** | HTTP handlers, auth, core task/wallet flows |
| **4** | WebSocket, RabbitMQ, Redis integration |
| **5** | Admin, tests, Swagger, deployment |
