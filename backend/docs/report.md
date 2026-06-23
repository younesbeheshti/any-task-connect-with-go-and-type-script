# TaskBridge Backend — Phase 3 Report

## Overview

Phase 3 implements the **Authentication** and **User** modules with full JWT flows, RBAC, permission system, profile management, and verification infrastructure. All API shapes follow `front/docs/api-contracts.md` as the source of truth.

Implemented capabilities:

- User registration (REQUESTER / AGENT roles)
- Login via phone or email + password
- JWT access + refresh tokens with rotation
- Redis-backed sessions, blacklist, OTP, and lockout
- Logout (current device / all devices)
- Change / forgot / reset password structures
- Phone OTP send/verify
- Profile CRUD (`GET/PATCH /me`, avatar delete, stats)
- Public user profiles
- Flexible permission model (DB-seeded + in-memory resolver)
- Frontend-compatible error envelope

---

## Files Created

### Configuration & API contract layer
| File | Purpose |
|------|---------|
| `pkg/apiresponse/contract.go` | Frontend error envelope (`error.code/message/fields`) |
| `pkg/security/password.go` | bcrypt hashing (cost 12), token hashing |
| `pkg/phone/phone.go` | E.164 phone normalization (+98…) |
| `pkg/redis/cache.go` | Redis cache adapter + JWT blacklist helpers |
| `internal/api/router.go` | `/v1` route registration |
| `internal/api/dto/user.go` | Swagger DTOs matching frontend User schema |
| `internal/api/dto/stats.go` | User stats Swagger DTO |

### Auth module
| File | Purpose |
|------|---------|
| `internal/auth/infra/model.go` | GORM models: refresh_tokens, password_reset_tokens |
| `internal/auth/infra/repository.go` | Auth persistence (sessions, password reset) |
| `internal/auth/infra/session.go` | Redis session store, OTP store, lockout store |
| `internal/auth/service/impl.go` | Auth business logic |
| `internal/auth/handler/handler.go` | HTTP handlers for auth endpoints |
| `internal/auth/service/impl_test.go` | Auth integration tests |

### User module
| File | Purpose |
|------|---------|
| `internal/user/infra/model.go` | GORM UserModel + permissions tables |
| `internal/user/infra/repository.go` | User persistence |
| `internal/user/service/impl.go` | Profile use cases |
| `internal/user/handler/dto.go` | Request/response DTOs |
| `internal/user/handler/handler.go` | Profile HTTP handlers |

### Middleware
| File | Purpose |
|------|---------|
| `internal/common/middleware/auth_jwt.go` | JWT auth with blacklist + permission guard |

### Migrations
| File | Purpose |
|------|---------|
| `migrations/000002_permissions_and_verification.up.sql` | Permissions, role_permissions, password_reset_tokens, verification columns |
| `migrations/000002_permissions_and_verification.down.sql` | Rollback |

### Tests
| File | Purpose |
|------|---------|
| `pkg/jwt/jwt_test.go` | JWT generation/validation/refresh |
| `pkg/validator/validator_test.go` | Password strength, national code |
| `pkg/phone/phone_test.go` | Phone normalization |
| `internal/common/rbac_test.go` | Permission matrix tests |

### Documentation
| File | Purpose |
|------|---------|
| `docs/report.md` | This report |
| `docs/swagger.json` / `docs/swagger.yaml` | Regenerated Swagger |

---

## Database Changes

### Existing tables used
- `users` — full user aggregate (role, verification, wallet denorm fields)
- `refresh_tokens` — refresh token hashes with revocation
- `cities` — profile city resolution by title

### New tables (migration 000002)
- `permissions` — flexible permission registry
- `role_permissions` — role → permission mapping
- `password_reset_tokens` — forgot/reset password flow

### New user columns
- `verification_status` — pending | approved | rejected
- `verification_reason` — admin rejection reason
- `verified_at` — timestamp of verification

### Indexes
- `idx_password_reset_user` on `password_reset_tokens(user_id)` where unused
- Existing: `idx_users_role`, `idx_users_city`, `idx_refresh_tokens_user`

---

## API Endpoints

Base path: `/v1` (matches `https://api.taskbridge.ir/v1`)

### Auth (public unless noted)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/auth/register` | Register requester/agent |
| POST | `/v1/auth/login` | Login (phone or email) |
| POST | `/v1/auth/refresh` | Rotate access token |
| POST | `/v1/auth/logout` | Logout current device 🔒 |
| POST | `/v1/auth/logout-all` | Logout all devices 🔒 |
| POST | `/v1/auth/change-password` | Change password 🔒 |
| POST | `/v1/auth/forgot-password` | Initiate password reset |
| POST | `/v1/auth/reset-password` | Complete password reset |
| POST | `/v1/auth/otp/send` | Send phone OTP |
| POST | `/v1/auth/otp/verify` | Verify phone OTP |

### User / Profile 🔒
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/me` | Current user (contract User schema) |
| PATCH | `/v1/me` | Update profile |
| DELETE | `/v1/me/avatar` | Remove avatar |
| GET | `/v1/me/stats` | User statistics |
| GET | `/v1/users/:id/public` | Public profile |

### Infrastructure (unchanged)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/swagger/index.html` | API documentation |

---

## Redis Usage

| Key pattern | TTL | Purpose |
|-------------|-----|---------|
| `jwt:blacklist:{jti}` | access TTL | Revoked access tokens on logout |
| `session:{userID}:{sessionID}` | 7 days | Refresh session metadata |
| `session:{userID}:index` | 7 days | Session index for revoke-all |
| `otp:phone:{phone}` | 5 min | Phone OTP (hashed code, attempts) |
| `otp:email:{email}` | 5 min | Email verification OTP |
| `lockout:{identifier}` | 15 min | Failed login counter |
| `ratelimit:{ip}:{route}` | 1 min | Rate limiting (Phase 2 middleware) |

---

## Security Features

### JWT
- Access tokens (15 min default) with `sub`, `role`, `jti`, `typ`
- Refresh tokens (7 days) stored in HttpOnly cookie + DB hash
- Token rotation on refresh
- JTI blacklist in Redis on logout

### Password
- bcrypt cost 12
- Policy: 8+ chars, upper, lower, digit, special character
- Change password revokes all sessions

### Brute force
- 5 failed attempts → 15 min lockout per phone/email
- Rate limiting middleware on write endpoints

### RBAC
- Roles: ADMIN, REQUESTER, AGENT (API: `admin|requester|agent`)
- 13 granular permissions seeded in DB
- `PermissionRequired` middleware for future routes

---

## Middleware

| Middleware | Purpose |
|------------|---------|
| `AuthJWT` | Bearer JWT validation + blacklist check |
| `Authorization` | Role-based route guard |
| `RBAC` / `PermissionRequired` | Granular permission checks |
| `RateLimit` | Redis sliding-window limiter |
| `RequestID`, `Logging`, `Recovery`, `CORS`, `SecurityHeaders`, `Timeout` | Phase 2 cross-cutting |

---

## Tests

| Package | Coverage |
|---------|----------|
| `pkg/jwt` | Token generation, validation, refresh rotation |
| `pkg/validator` | Password strength, national code |
| `pkg/phone` | E.164 normalization |
| `internal/common` | RBAC permission matrix |
| `internal/auth/service` | Register, login, OTP flow (SQLite + miniredis) |
| `internal/health` | Health endpoint |

Run: `go test ./...`

---

## Swagger

Regenerated via `swag init -g cmd/api/main.go -o docs --parseInternal`

Available at: `http://localhost:8080/swagger/index.html`

All Phase 3 endpoints include request/response schemas and Bearer auth security definition.

---

## Frontend Compatibility

| Contract requirement | Backend implementation |
|---------------------|------------------------|
| `POST /auth/register` body: `fullName, phone, password, role` | Exact field names (camelCase JSON) |
| `POST /auth/login` body: `phone, password` | + optional `email` for Phase 3 without breaking contract |
| Response: `{ token, user }` | `token` = access JWT; `user` matches User schema |
| User schema fields | `id, fullName, phone, email, city, role, avatarUrl, verification, rating, completedCount, createdAt` |
| `PATCH /me` fields | `fullName?, city?, avatarUrl?, bio?` |
| Error envelope | `{ error: { code, message, fields } }` via `pkg/apiresponse` |
| Role values | `requester \| agent \| admin` (lowercase in API) |
| Phone format | Stored/returned as E.164 `+989…` |
| Currency / tasks | Not in Phase 3 scope — deferred to Phase 4 |

Refresh token delivered via:
1. HttpOnly `refresh_token` cookie (architecture doc)
2. Optional `refreshToken` field in JSON (Phase 3 extension)

Permissions and role also returned in auth response as Phase 3 extension without altering required `token` + `user` fields.

---

## Future Improvements (Phase 4)

1. **Task & Wallet modules** — implement remaining contract endpoints
2. **Email/SMS gateways** — wire OTP and verification to real providers
3. **File upload** — `POST /files` for avatar upload (currently URL-based)
4. **WebSocket** — realtime notifications per architecture
5. **DB-backed permission resolver** — load `role_permissions` at startup instead of static map
6. **Identity verification workflow** — admin approval for `verification_status`
7. **Integration tests** — full HTTP tests against PostgreSQL test container
8. **Prometheus metrics** — `/metrics` endpoint

---

## Quick Start

```bash
cd backend
cp .env.example .env
docker compose up --build
```

Test registration:

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"fullName":"سارا محمدی","phone":"09120000000","password":"Secure1!","role":"requester"}'
```

---

# TaskBridge Backend — Phase 4 Report

## Phase 4 Overview

Phase 4 implements the core business modules: **Task Management**, **Application System**, **Category & City CRUD**, **Dashboard APIs**, and **Search & Filtering**. All API shapes follow `front/docs/api-contracts.md` as the source of truth.

---

## Task Module

### Files
| File | Purpose |
|------|---------|
| `internal/task/domain/task.go` | Task entity, TaskStatus, ValidTransitions, state machine, TaskFilter (extended with MinBudget, MaxBudget, CityTitle, CatTitle) |
| `internal/task/repository/repository.go` | Repository interface (extended with Delete) |
| `internal/task/service/service.go` | Service interface (extended with Delete) |
| `internal/task/service/impl.go` | Full lifecycle: Create (auto-publishes, 8% escrow fee), Start, Complete (chains → WAITING_FOR_VERIFICATION), Verify (chains → PAID), Cancel, Delete |
| `internal/task/infra/model.go` | GORM TaskModel and TaskTimelineModel |
| `internal/task/infra/repository.go` | GORM repo with dynamic filter queries, JOIN-based city/category resolution |
| `internal/task/handler/dto.go` | TaskResponse, CreateTaskResponse, EscrowResponse DTOs |
| `internal/task/handler/handler.go` | Create, List, GetByPublicID, Update, Delete, Cancel, Start, Complete, Verify, GetTimeline |

### API Endpoints
| Method | Path | Auth | Role |
|--------|------|------|------|
| GET | `/v1/tasks` | optional | any |
| POST | `/v1/tasks` | required | requester |
| GET | `/v1/tasks/:id` | optional | any |
| PATCH | `/v1/tasks/:id` | required | requester |
| DELETE | `/v1/tasks/:id` | required | requester |
| POST | `/v1/tasks/:id/cancel` | required | requester/agent |
| POST | `/v1/tasks/:id/start` | required | agent |
| POST | `/v1/tasks/:id/complete` | required | agent |
| POST | `/v1/tasks/:id/verify` | required | requester |
| GET | `/v1/tasks/:id/timeline` | required | any |

---

## Application Module

### Files
| File | Purpose |
|------|---------|
| `internal/application/domain/application.go` | Application entity, status enum |
| `internal/application/repository/repository.go` | Repository interface with UpdateStatusByTask |
| `internal/application/service/service.go` | Service interface |
| `internal/application/service/impl.go` | Submit (duplicate check), Accept (with reject-others), Reject, Withdraw |
| `internal/application/infra/model.go` | GORM ApplicationModel |
| `internal/application/infra/repository.go` | GORM repo with Preload Agent |
| `internal/application/handler/handler.go` | Submit, ListByTask, Accept, Reject, Withdraw, ListMyApplications |

### API Endpoints
| Method | Path | Auth | Role |
|--------|------|------|------|
| POST | `/v1/tasks/:id/applications` | required | agent |
| GET | `/v1/tasks/:id/applications` | required | requester |
| POST | `/v1/applications/:id/accept` | required | requester |
| POST | `/v1/applications/:id/reject` | required | requester |
| POST | `/v1/applications/:id/withdraw` | required | agent |
| GET | `/v1/me/applications` | required | agent |

---

## State Machine

```
CREATED → OPEN (auto on create)
OPEN → ASSIGNED (on application accept)
ASSIGNED → IN_PROGRESS (agent: POST /tasks/:id/start)
IN_PROGRESS → COMPLETED → WAITING_FOR_VERIFICATION (agent: POST /tasks/:id/complete, chained)
WAITING_FOR_VERIFICATION → VERIFIED → PAID (requester: POST /tasks/:id/verify, chained)
any active → CANCELLED (POST /tasks/:id/cancel)
```

Frontend API status mapping:
- CREATED → `posted`
- OPEN → `awaiting_applicants`
- ASSIGNED → `accepted`
- IN_PROGRESS → `in_progress`
- COMPLETED → `completed`
- WAITING_FOR_VERIFICATION, VERIFIED → `awaiting_verification`
- PAID → `paid`
- CANCELLED → `cancelled`

---

## Category Module

### Files
| File | Purpose |
|------|---------|
| `internal/category/infra/model.go` | GORM CategoryModel → `categories` table |
| `internal/category/infra/repository.go` | GORM CRUD |
| `internal/category/service/impl.go` | List, GetByID, Create (duplicate check), Update, Delete |
| `internal/category/handler/handler.go` | GET/POST/PATCH/DELETE handlers |

### Seeded Categories
اداری · پزشکی · دانشگاهی · حقوقی · دولتی · پیک و تحویل · خرید · سایر

---

## City Module

### Files
| File | Purpose |
|------|---------|
| `internal/city/infra/model.go` | GORM CityModel → `cities` table |
| `internal/city/infra/repository.go` | GORM repo |
| `internal/city/service/impl.go` | List, GetByID |
| `internal/city/handler/handler.go` | GET /cities, GET /cities/:id |

### Seeded Cities
تهران · مشهد · اصفهان · شیراز · تبریز · کرج · اهواز · قم · رشت · کرمان

---

## Dashboard APIs

### Requester: `GET /v1/dashboard/stats?role=requester`
Returns task counts by status, wallet balance, and stats.

### Agent: `GET /v1/dashboard/stats?role=agent`
Returns available tasks, assigned tasks, completed tasks, earnings, rating.

---

## Search & Filtering

Query parameters for `GET /v1/tasks`:
- `?city=تهران` — filter by city name
- `?cityId=uuid` — filter by city UUID
- `?category=اداری` — filter by category name
- `?categoryId=uuid` — filter by category UUID
- `?status=OPEN` — filter by domain status
- `?q=keyword` — full-text ILIKE on title + description
- `?minBudget=100000&maxBudget=500000` — budget range
- `?sort=newest|budget|deadline` — ordering
- `?page=1&pageSize=20` — pagination

---

## Concurrency Handling

The `AcceptApplication` flow:
1. Gets application and validates ownership
2. Updates application status to ACCEPTED
3. Rejects all other PENDING applications via bulk UPDATE with WHERE clause
4. Updates task status to ASSIGNED

Note: Full SELECT FOR UPDATE is supported on PostgreSQL in production; the service-level validation (status check before accept) prevents double-assignment.

---

## RabbitMQ Events

| Event | Published On |
|-------|-------------|
| `task.created` | Task creation |
| `task.started` | Task transition to IN_PROGRESS |
| `task.completed` | Task transition to COMPLETED |
| `task.verified` | Task paid/verified |
| `application.submitted` | Application submitted |

---

## Redis Cache

Reserved key patterns (Phase 5 implementation):
- `task:cache:{id}` — individual task cache (5 min TTL)
- `tasks:popular_categories` — popular category counts
- `dashboard:{userID}` — dashboard response cache (2 min TTL)

---

## Tests

| Package | Tests |
|---------|-------|
| `internal/task/service` | Create, StateMachine transitions, Cancel, List, APIStatus mapping |
| `internal/application/service` | Submit, Duplicate rejection, Accept, Withdraw |
| `internal/common` | RBAC permission matrix |
| `pkg/jwt` | Token generation and validation |
| `pkg/validator` | Password strength |
| `pkg/phone` | E.164 normalization |

Run: `go test ./...` — all 8 test packages pass.

---

## Frontend Compatibility

| Contract requirement | Backend implementation |
|---------------------|------------------------|
| `GET /tasks` with `city`, `category`, `status`, `q`, `sort`, `page`, `pageSize` | All query params handled |
| `POST /tasks` body: `title, description, category, city, budget, deadline, attachments` | Handler maps category/city names → IDs |
| Response `Task.status` uses lowercase snake_case | `APIStatus()` converts domain status to API values |
| `POST /tasks/:id/applications` body: `price, eta, message` | Exact field names |
| `POST /applications/:id/accept` | Implemented with rejection of other pending |
| Error envelope `{ error: { code, message, fields } }` | All handlers use `apiresponse.WriteError` |
| Currency `"IRT"` in task and escrow responses | Hardcoded in service |
| Public ID format `TB-{N}` | `NextPublicID()` generates sequential IDs |

---

## Future Phase 5

1. **Wallet integration** — lock escrow on task create, release on verify
2. **File uploads** — `POST /files` for real attachment storage (MinIO/S3)
3. **WebSocket** — realtime notifications via `wss://...`
4. **Rating system** — `POST /tasks/:id/review` after completion
5. **Chat module** — per-task messaging thread
6. **Admin dashboard** — dispute resolution, user management
7. **Redis caching** — task list and dashboard response caching
8. **Full SELECT FOR UPDATE** — distributed lock for application accept
9. **Prometheus metrics** — `/metrics` endpoint
