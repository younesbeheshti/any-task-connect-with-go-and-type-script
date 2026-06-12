# TaskBridge — API Contracts & JSON Schemas

REST-style, all responses JSON. Currency: **Toman** (integer, no decimals). Dates: ISO-8601 for transport; UI renders Jalali.

Base URL: `https://api.taskbridge.ir/v1`

Auth: `Authorization: Bearer <jwt>` on every protected route. Role claim is one of `requester | agent | admin`.

---

## 1. Auth

### `POST /auth/register`
Request:
```json
{
  "fullName": "سارا محمدی",
  "phone": "+989120000000",
  "password": "string(min 8)",
  "role": "requester | agent"
}
```
Response 201:
```json
{ "token": "jwt", "user": { "$ref": "#/User" } }
```

### `POST /auth/login`
```json
{ "phone": "+989120000000", "password": "string" }
```
Response 200: same as register.

### `POST /auth/otp/send` / `POST /auth/otp/verify`
```json
{ "phone": "+989120000000" }
{ "phone": "+989120000000", "code": "123456" }
```

---

## 2. Users & Profile

### `GET /me`
Response: `User`.

### `PATCH /me`
```json
{ "fullName": "string?", "city": "string?", "avatarUrl": "string?", "bio": "string?" }
```

### `GET /users/:id/public`
Public profile for ratings & history.

### Schema `User`
```json
{
  "id": "uuid",
  "fullName": "string",
  "phone": "+98...",
  "email": "string|null",
  "city": "تهران|مشهد|...",
  "role": "requester|agent|admin",
  "avatarUrl": "string|null",
  "verification": {
    "phone": true,
    "email": false,
    "nationalId": true,
    "level": "none|basic|full"
  },
  "rating": 4.9,
  "completedCount": 27,
  "createdAt": "2026-06-01T08:00:00Z"
}
```

---

## 3. Tasks

### `GET /tasks`
Query: `city, category, status, sort=newest|budget|deadline, q, page, pageSize`.

Response:
```json
{ "items": [{ "$ref": "#/Task" }], "page": 1, "pageSize": 20, "total": 124 }
```

### `POST /tasks` (requester)
```json
{
  "title": "string(min 6, max 120)",
  "description": "string(min 20, max 2000)",
  "category": "اداری|پزشکی|دانشگاهی|حقوقی|دولتی|پیک و تحویل|خرید|سایر",
  "city": "تهران|...",
  "budget": 750000,
  "deadline": "2026-06-15",
  "attachments": ["file-id-1", "file-id-2"]
}
```
Response 201: `Task` + escrow info:
```json
{
  "task": { "$ref": "#/Task" },
  "escrow": { "fee": 60000, "held": 810000, "currency": "IRT" }
}
```

### `GET /tasks/:id`
Returns `Task` with `applicants[]`, `timeline[]`, `assignedAgent?`.

### `PATCH /tasks/:id`
Requester only, while status `posted | awaiting_applicants`.

### `POST /tasks/:id/cancel`
Returns refund result.

### `POST /tasks/:id/verify` (requester)
Moves status `awaiting_verification → paid`, releases escrow.

### `POST /tasks/:id/dispute`
```json
{ "reason": "string", "evidence": ["file-id"] }
```

### Schema `Task`
```json
{
  "id": "TB-1042",
  "title": "string",
  "description": "string",
  "category": "string",
  "city": "string",
  "budget": 750000,
  "currency": "IRT",
  "deadline": "2026-06-15",
  "status": "posted|awaiting_applicants|accepted|in_progress|completed|awaiting_verification|paid|cancelled",
  "postedBy": { "id": "uuid", "fullName": "string", "avatarUrl": "string|null" },
  "assignedAgentId": "uuid|null",
  "applicantsCount": 7,
  "attachments": [{ "id": "uuid", "name": "string", "size": 12345, "url": "string" }],
  "createdAt": "iso",
  "updatedAt": "iso"
}
```

---

## 4. Applications

### `POST /tasks/:id/applications` (agent)
```json
{ "price": 800000, "eta": "today|tomorrow|YYYY-MM-DD", "message": "string" }
```

### `GET /tasks/:id/applications` (requester)
Returns `Application[]`.

### `POST /applications/:id/accept` (requester)
Moves task to `accepted`, assigns agent.

### Schema `Application`
```json
{
  "id": "uuid",
  "taskId": "TB-1042",
  "agent": { "$ref": "#/User" },
  "price": 800000,
  "eta": "string",
  "message": "string",
  "createdAt": "iso",
  "status": "pending|accepted|rejected|withdrawn"
}
```

---

## 5. Wallet & Transactions

### `GET /wallet`
```json
{
  "availableBalance": 12400000,
  "lockedBalance": 3650000,
  "currency": "IRT",
  "bankCards": [{ "id": "uuid", "last4": "1234", "brand": "ملت" }]
}
```

### `POST /wallet/topup`
```json
{ "amount": 5000000, "cardId": "uuid|null", "returnUrl": "https://..." }
```
Response: `{ "paymentUrl": "https://pay.ir/..." }`.

### `POST /wallet/withdraw`
```json
{ "amount": 2000000, "iban": "IR123456789012345678901234" }
```

### `GET /transactions`
Query: `type, dateFrom, dateTo, page`.

### Schema `Transaction`
```json
{
  "id": "tx-901",
  "date": "2026-06-09T10:00:00Z",
  "description": "string",
  "amount": -850000,
  "type": "escrow_in|release|topup|refund|withdraw|fee",
  "status": "pending|completed|failed|locked|released",
  "relatedTaskId": "TB-1042|null"
}
```

---

## 6. Messaging

### `GET /chats`
Returns `Chat[]` (one per task or per pair).

### `GET /chats/:id/messages?before=cursor&limit=50`

### `POST /chats/:id/messages`
```json
{ "text": "string?", "attachments": ["file-id"]?, "voiceUrl": "string?" }
```

### Schema `Message`
```json
{
  "id": "uuid",
  "chatId": "uuid",
  "from": "uuid",
  "text": "string|null",
  "attachments": [{ "id": "uuid", "url": "string", "mime": "string" }],
  "voiceUrl": "string|null",
  "readAt": "iso|null",
  "createdAt": "iso"
}
```

WebSocket: `wss://api.taskbridge.ir/ws` — events `message.created`, `message.read`, `task.status_changed`, `notification.created`.

---

## 7. Notifications

### `GET /notifications?unread=true&page=1`
### `POST /notifications/read-all`
### `POST /notifications/:id/read`

### Schema `Notification`
```json
{
  "id": "uuid",
  "type": "task_assigned|task_completed|payment|application|message|dispute|system",
  "title": "string",
  "desc": "string",
  "unread": true,
  "createdAt": "iso",
  "deepLink": "/app/tasks/TB-1042"
}
```

---

## 8. Admin

All require `role=admin`.

### `GET /admin/metrics`
```json
{
  "users": { "total": 18402, "growth7d": 0.052 },
  "tasks": { "active": 847, "growth7d": 0.12 },
  "revenue": { "monthToDate": 9214000000, "growth7d": 0.084 },
  "disputes": { "open": 6, "delta7d": -2 }
}
```

### `GET /admin/users` — paginated user search/list.
### `PATCH /admin/users/:id` — `{ status: "active|suspended", roles: [...] }`.
### `GET /admin/transactions` — global ledger.
### `GET /admin/reports/revenue?range=30d|90d|year`
### `POST /admin/disputes/:id/resolve` — `{ outcome: "refund_requester|release_agent|split", note: "string" }`.

---

## 9. Files / Uploads

### `POST /files`
Multipart `file`. Response:
```json
{ "id": "uuid", "url": "https://cdn.taskbridge.ir/...", "mime": "image/jpeg", "size": 12345 }
```

---

## 10. Status State Machine

```
posted
  └─> awaiting_applicants ──> accepted ──> in_progress ──> completed
                                                              │
                                                              ▼
                                                awaiting_verification
                                                              │
                                                              ▼
                                                            paid
any active state ──> cancelled (with refund rules)
```

Allowed transitions enforced server-side; UI mirrors them.

---

## 11. Error envelope
```json
{
  "error": {
    "code": "VALIDATION_FAILED|FORBIDDEN|NOT_FOUND|RATE_LIMITED|PAYMENT_FAILED|...",
    "message": "Human-readable Persian message",
    "fields": { "budget": "مبلغ باید مثبت باشد" }
  }
}
```

HTTP status codes: 200/201/204 on success; 400/401/403/404/409/422/429/500 on errors.
