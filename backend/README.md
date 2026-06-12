# TaskBridge Backend — Phase 1

Production backend for the TaskBridge marketplace (escrow-based task platform).

## Phase 1 Deliverables (Design Foundation)

This phase establishes the architecture and contracts **without** full implementation.

| Artifact | Location |
|----------|----------|
| Architecture | [docs/architecture.md](docs/architecture.md) |
| ER Diagram | [docs/erd.md](docs/erd.md) |
| State Machine | [docs/state-machine.md](docs/state-machine.md) |
| Dependency Graph | [docs/dependency-graph.md](docs/dependency-graph.md) |
| Sequence Diagrams | [docs/sequence-diagrams.md](docs/sequence-diagrams.md) |
| Database Schema | [migrations/000001_init_schema.up.sql](migrations/000001_init_schema.up.sql) |
| Domain Models | `internal/*/domain/` |
| Repository Interfaces | `internal/*/repository/` |
| Service Interfaces | `internal/*/service/` |

## Project Structure

```
backend/
├── cmd/api/                 # Application entrypoint (Phase 2)
├── internal/
│   ├── admin/               # Admin analytics & moderation
│   ├── application/         # Agent task applications
│   ├── auth/                # JWT, OTP, password reset
│   ├── category/            # Task categories
│   ├── chat/                # Task-scoped messaging
│   ├── city/                # Geographic locations
│   ├── common/              # Shared types, RBAC, events, cache
│   ├── notification/        # In-app notifications
│   ├── payment/             # Escrow & transactions
│   ├── rating/              # Reviews
│   ├── task/                # Task lifecycle & state machine
│   ├── user/                # User profiles
│   └── wallet/              # Wallet balances
├── pkg/                     # Reusable infrastructure (Phase 2)
│   ├── configs/
│   ├── jwt/
│   ├── redis/
│   ├── rabbitmq/
│   └── websocket/
├── migrations/              # golang-migrate SQL
├── docs/                    # Architecture & diagrams
├── scripts/                 # Dev/ops scripts
└── docker/                  # Docker Compose (Phase 2)
```

## Tech Stack

Go 1.24 · Gin · PostgreSQL · GORM · Redis · RabbitMQ · JWT · Viper · Zap · validator · golang-migrate · testify · Docker

## Frontend API Alignment

Contracts defined in [`../front/docs/api-contracts.md`](../front/docs/api-contracts.md). Domain states map to API statuses via `TaskStatus.APIStatus()`.

## Next Phase (awaiting approval)

1. Infrastructure adapters (GORM repositories, Redis, RabbitMQ)
2. Service implementations with transaction boundaries
3. HTTP handlers + middleware + `/api/v1` routes
4. WebSocket hub
5. Tests, Swagger, Docker Compose

## Approval Checklist

- [ ] Architecture layers and bounded contexts
- [ ] ER diagram and PostgreSQL schema
- [ ] Task state machine transitions
- [ ] Domain models match entities
- [ ] Repository/service interface coverage
- [ ] Dependency graph and event flow

Reply **approved** (or with changes) to proceed to Phase 2 implementation.
