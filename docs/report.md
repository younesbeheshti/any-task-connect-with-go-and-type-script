# TaskBridge — Platform Design Report

**Date:** 2026-06-24  
**Scope:** Full role-based platform audit. Covers UI structure, backend endpoints, permission matrix, gap analysis, and future-phase recommendations.

---

## 1. Role-Based UI Structure

### 1.1 Requester

The requester creates tasks, posts them publicly, receives applications, accepts one agent, monitors task progress, and pays on completion.

**Navigation (sidebar):**

| Route | Label | Icon |
|---|---|---|
| `/app` | داشبورد | LayoutDashboard |
| `/app/tasks/new` | ثبت درخواست جدید | PlusCircle |
| `/app/tasks` | درخواست‌های من | ClipboardList |
| `/app/in-progress` | در حال انجام | Activity |
| `/app/completed` | تکمیل‌شده‌ها | CheckCircle2 |
| `/app/wallet` | کیف پول | Wallet |
| `/app/reviews` | نظرات و امتیازها | Star |
| `/app/chat` | گفت‌وگوها | MessagesSquare |
| `/app/notifications` | اعلان‌ها | Bell |
| `/app/profile` | پروفایل | User |

**Key flows:**
1. Create task → `POST /v1/tasks`
2. View applications → `GET /v1/tasks/:id/applications`
3. Accept/reject agent → `POST /v1/applications/:id/accept|reject`
4. Monitor progress → Task detail page with timeline (`GET /v1/tasks/:id/timeline`)
5. Verify completion → `POST /v1/tasks/:id/verify`
6. Leave review → `POST /v1/tasks/:id/reviews`

### 1.2 Agent

The agent browses open tasks, submits applications, performs assigned work, marks complete, and withdraws earnings.

**Navigation (sidebar):**

| Route | Label | Icon |
|---|---|---|
| `/app/agent` | داشبورد | LayoutDashboard |
| `/app/tasks` | فرصت‌های جدید | Sparkles |
| `/app/applications` | درخواست‌های من | FileText |
| `/app/accepted` | پذیرفته‌شده‌ها | Briefcase |
| `/app/completed` | تکمیل‌شده‌ها | CheckCircle2 |
| `/app/earnings` | درآمدها | BarChart3 |
| `/app/wallet` | کیف پول | Wallet |
| `/app/reviews` | نظرات و امتیازها | Star |
| `/app/chat` | گفت‌وگوها | MessagesSquare |
| `/app/notifications` | اعلان‌ها | Bell |
| `/app/profile` | پروفایل | User |

**Key flows:**
1. Browse open tasks → `GET /v1/tasks?status=OPEN`
2. Apply to task → `POST /v1/tasks/:id/applications`
3. Withdraw application → `POST /v1/applications/:id/withdraw`
4. Start assigned task → `POST /v1/tasks/:id/start`
5. Mark task complete → `POST /v1/tasks/:id/complete`
6. Withdraw earnings → `POST /v1/wallet/withdraw`

### 1.3 Admin

The admin moderates platform health — user management, financial oversight, task governance.

**Navigation (sidebar):**

| Route | Label | Icon |
|---|---|---|
| `/app/admin` | داشبورد مدیریت | LayoutDashboard |
| `/app/admin/users` | کاربران | Users |
| `/app/tasks` | درخواست‌ها | ClipboardList |
| `/app/admin/finance` | مالی و تراکنش‌ها | Wallet |
| `/app/admin/reports` | گزارش‌ها | FileBarChart |
| `/app/notifications` | اعلان‌ها | Bell |
| `/app/profile` | پروفایل | User |

---

## 2. Page Hierarchy

```
/                          → Landing / marketing
/login                     → Login form
/register                  → Registration form

/app                       → Authenticated shell (role-aware sidebar)
  /app                     → Requester dashboard
  /app/agent               → Agent dashboard
  /app/admin               → Admin dashboard

  /app/tasks               → Task list (browse for agents, own tasks for requesters)
  /app/tasks/new           → Create task form
  /app/tasks/:id           → Task detail
  /app/tasks/:id/applications → Requester reviews applicants

  /app/in-progress         → Requester's ASSIGNED/IN_PROGRESS tasks
  /app/accepted            → Agent's ASSIGNED/IN_PROGRESS tasks
  /app/completed           → Completed tasks (both roles)

  /app/applications        → Agent's submitted applications
  /app/wallet              → Wallet + top-up + history
  /app/earnings            → Agent earnings breakdown
  /app/reviews             → Reviews received by current user
  /app/chat                → Chat inbox + task-scoped conversations
  /app/notifications       → Notification list (filterable)
  /app/profile             → Edit profile, avatar, skills

  /app/admin/users         → User table with search, suspend, activate
  /app/admin/finance       → Withdraw requests + transaction log
  /app/admin/reports       → KPI cards, exportable metrics
```

---

## 3. User Journeys

### Requester journey
```
Register → Verify Email → Top Up Wallet → Create Task →
  Receive Applications → Chat with Applicant → Accept Agent →
  Monitor Progress → Mark for Verification → Verify Task →
  Leave Review → Withdraw Balance
```

### Agent journey
```
Register → Complete Profile → Browse Open Tasks → Apply →
  Chat with Requester → Get Accepted → Start Work →
  Mark Complete → Await Verification → Receive Payment →
  Withdraw Earnings → Build Rating
```

### Admin journey
```
Login → Review Dashboard KPIs → Investigate Disputes →
  Manage Users (suspend/activate) → Approve/Reject Withdrawals →
  Monitor Revenue → Export Reports
```

---

## 4. Permission Matrix

| Permission | REQUESTER | AGENT | ADMIN |
|---|---|---|---|
| View public tasks | ✓ | ✓ | ✓ |
| Create task | ✓ | — | — |
| Update own task | ✓ | — | ✓ |
| Delete own task | ✓ | — | ✓ |
| Submit application | — | ✓ | — |
| Withdraw application | — | ✓ | — |
| Accept/Reject application | ✓ | — | ✓ |
| Start task | — | ✓ | — |
| Complete task | — | ✓ | — |
| Verify task | ✓ | — | ✓ |
| Cancel task | ✓ | — | ✓ |
| View own wallet | ✓ | ✓ | — |
| Top up wallet | ✓ | ✓ | — |
| Withdraw from wallet | ✓ | ✓ | — |
| Send/receive messages | ✓ | ✓ | — |
| Create review | ✓ | ✓ | — |
| View admin metrics | — | — | ✓ |
| Manage users | — | — | ✓ |
| Approve/reject withdrawals | — | — | ✓ |
| Manage categories/cities | — | — | ✓ |

---

## 5. Backend Endpoint List

### Auth

| Method | Path | Description |
|---|---|---|
| POST | `/v1/auth/register` | Register new user |
| POST | `/v1/auth/login` | Login, return JWT pair |
| POST | `/v1/auth/refresh` | Refresh access token |
| POST | `/v1/auth/logout` | Blacklist JTI |

### User

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/v1/me` | ✓ | Own profile |
| PATCH | `/v1/me` | ✓ | Update profile |
| GET | `/v1/users/:id` | public | Public profile |
| GET | `/v1/users/:id/reviews` | public | Reviews for user |

### Task

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/v1/tasks` | public | List tasks (filterable) |
| POST | `/v1/tasks` | ✓ | Create task |
| GET | `/v1/tasks/:id` | public | Task detail |
| PATCH | `/v1/tasks/:id` | ✓ | Update task |
| DELETE | `/v1/tasks/:id` | ✓ | Delete task |
| GET | `/v1/tasks/:id/timeline` | ✓ | Task status history |
| POST | `/v1/tasks/:id/cancel` | ✓ | Cancel task |
| POST | `/v1/tasks/:id/start` | ✓ | Agent starts work |
| POST | `/v1/tasks/:id/complete` | ✓ | Agent marks complete |
| POST | `/v1/tasks/:id/verify` | ✓ | Requester verifies |

### Application

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/v1/tasks/:id/applications` | ✓ | Submit application |
| GET | `/v1/tasks/:id/applications` | ✓ | Applications for task |
| GET | `/v1/applications/:id` | ✓ | Single application |
| GET | `/v1/me/applications` | ✓ | My applications (agent) |
| POST | `/v1/applications/:id/accept` | ✓ | Requester accepts |
| POST | `/v1/applications/:id/reject` | ✓ | Requester rejects |
| POST | `/v1/applications/:id/withdraw` | ✓ | Agent withdraws |

### Dashboard

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/v1/dashboard/stats` | ✓ | User stats (role-aware) |
| GET | `/v1/admin/dashboard/admin-stats` | ADMIN | Platform aggregate stats |

### Wallet

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/v1/wallet` | ✓ | Wallet balances |
| GET | `/v1/wallet/history` | ✓ | Transaction history |
| GET | `/v1/wallet/statistics` | ✓ | Earnings stats |
| POST | `/v1/wallet/topup` | ✓ | Add balance |
| POST | `/v1/wallet/withdraw` | ✓ | Create withdraw request |
| GET | `/v1/withdraws` | ✓ | My withdraw requests |
| GET | `/v1/withdraws/:id` | ✓ | Single withdraw |

### Transaction

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/v1/transactions` | ✓ | Own transaction log |
| GET | `/v1/transactions/:id` | ✓ | Single transaction |

### Chat

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/v1/chats` | ✓ | Chat list |
| GET | `/v1/tasks/:id/messages` | ✓ | Messages (cursor-based) |
| POST | `/v1/tasks/:id/messages` | ✓ | Send message |
| POST | `/v1/tasks/:id/messages/read` | ✓ | Mark task messages read |

### Notification

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/v1/notifications` | ✓ | List notifications |
| PATCH | `/v1/notifications/:id/read` | ✓ | Mark one read |
| POST | `/v1/notifications/read-all` | ✓ | Mark all read |

### Rating

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/v1/tasks/:id/reviews` | ✓ | Submit review |
| GET | `/v1/tasks/:id/reviews` | public | Reviews for task |
| GET | `/v1/users/:id/reviews` | public | Reviews for user |

### Category / City

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/v1/categories` | public | List categories |
| GET | `/v1/categories/:id` | public | Single category |
| POST | `/v1/admin/categories` | ADMIN | Create |
| PATCH | `/v1/admin/categories/:id` | ADMIN | Update |
| DELETE | `/v1/admin/categories/:id` | ADMIN | Delete |
| GET | `/v1/cities` | public | List cities |
| GET | `/v1/cities/:id` | public | Single city |
| POST | `/v1/admin/cities` | ADMIN | Create |
| PATCH | `/v1/admin/cities/:id` | ADMIN | Update |
| DELETE | `/v1/admin/cities/:id` | ADMIN | Delete |

### Admin

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/v1/admin/metrics` | ADMIN | Platform KPIs |
| GET | `/v1/admin/users` | ADMIN | User search/list |
| GET | `/v1/admin/users/:id` | ADMIN | User detail |
| PATCH | `/v1/admin/users/:id` | ADMIN | Update user |
| POST | `/v1/admin/users/:id/suspend` | ADMIN | Suspend user |
| POST | `/v1/admin/users/:id/activate` | ADMIN | Activate user |
| GET | `/v1/admin/revenue` | ADMIN | Revenue totals |
| GET | `/v1/admin/revenue/statistics` | ADMIN | Revenue stats |
| GET | `/v1/admin/revenue/daily` | ADMIN | Daily breakdown |
| GET | `/v1/admin/revenue/monthly` | ADMIN | Monthly breakdown |
| GET | `/v1/admin/financial/overview` | ADMIN | Financial overview |
| GET | `/v1/admin/financial/transactions` | ADMIN | All transactions |
| GET | `/v1/admin/financial/withdraws` | ADMIN | All withdraw requests |
| POST | `/v1/admin/withdraws/:id/approve` | ADMIN | Approve withdraw |
| POST | `/v1/admin/withdraws/:id/reject` | ADMIN | Reject withdraw |

---

## 6. Database Entities

### Implemented

| Table | Context | Notes |
|---|---|---|
| `users` | auth/user | Roles: ADMIN, REQUESTER, AGENT |
| `tasks` | task | Full state machine |
| `task_timeline` | task | Status history per task |
| `applications` | application | Agent applications |
| `wallets` | wallet | Per-user, auto-created on first access |
| `transactions` | payment | Double-entry ledger rows |
| `withdraw_requests` | withdraw | Agent bank withdrawal requests |
| `revenue_entries` | revenue | Platform fee records |
| `chat_messages` | chat | Task-scoped messages |
| `notifications` | notification | Per-user push-style notifications |
| `reviews` | rating | Ratings 1–5 with optional comment |
| `categories` | category | Task taxonomy |
| `cities` | city | Geographic filter |

### Missing / Planned

| Entity | Purpose | Phase |
|---|---|---|
| `disputes` | Formal escrow dispute tickets | Phase 6 |
| `user_skills` | Agent skill tags for matching | Phase 6 |
| `audit_logs` | Admin action log for compliance | Phase 6 |
| `platform_settings` | Configurable fee %, limits | Phase 6 |
| `email_verifications` | Email OTP tokens | Phase 6 |
| `notification_preferences` | Per-user notification opt-ins | Phase 6 |

---

## 7. Missing APIs / Services / Events

### APIs not yet implemented

| Endpoint | Priority | Notes |
|---|---|---|
| `POST /v1/auth/verify-email` | High | Email confirmation flow |
| `POST /v1/auth/forgot-password` | High | Password reset |
| `GET /v1/users/:id` | Medium | Public agent profile card |
| `PATCH /v1/me/avatar` | Medium | Avatar upload (multipart) |
| `GET /v1/tasks/:id/applicant-count` | Low | Real-time badge for task detail |
| `POST /v1/admin/disputes/:id/resolve` | Medium | Currently returns 501 |
| `GET /v1/admin/audit-logs` | Medium | Admin action history |
| `POST /v1/admin/notifications/broadcast` | Low | Platform-wide push |

### Events not emitted (RabbitMQ stubs)

| Event | Trigger | Consumer |
|---|---|---|
| `task.created` | POST /tasks | Notify nearby agents |
| `application.accepted` | Accept application | Notify agent + lock funds |
| `task.completed` | POST /tasks/:id/complete | Notify requester |
| `task.verified` | POST /tasks/:id/verify | Trigger escrow release |
| `withdraw.approved` | Admin approves | Trigger bank transfer stub |
| `review.created` | POST /tasks/:id/reviews | Update user avg rating |

### Notification triggers not yet wired

Backend creates notifications only manually today. The following should be auto-created by event consumers:

- Task accepted (agent notified)
- New application on my task (requester notified)
- Task completed (requester notified)
- Withdraw approved/rejected (agent notified)
- New message received (recipient notified)

---

## 8. Feature Gap Analysis

### Authentication

| Feature | Status |
|---|---|
| Login / Register | ✓ Implemented |
| JWT access + refresh | ✓ Implemented |
| Logout (JTI blacklist) | ✓ Implemented |
| Email verification | ✗ Missing |
| Password reset | ✗ Missing |
| OAuth (Google) | ✗ Not planned |

### Task Lifecycle

| Feature | Status |
|---|---|
| Full state machine | ✓ Implemented |
| Task timeline | ✓ Implemented |
| Requester cancels | ✓ Implemented |
| Escrow lock on accept | ✓ Implemented |
| Escrow release on verify | ✓ Implemented |
| Platform fee deduction | ✓ Implemented (revenue_entries) |
| Task expiry / auto-cancel | ✗ Missing (cron job needed) |

### Chat

| Feature | Status |
|---|---|
| Task-scoped messages | ✓ Implemented |
| Cursor-based pagination | ✓ Implemented |
| Mark messages read | ✓ Implemented |
| File/image attachments | ✗ UI removed; backend supports JSON blob |
| Real-time WebSocket | ✗ Missing (currently polling only) |

### Notifications

| Feature | Status |
|---|---|
| Create / list / mark read | ✓ Implemented |
| Type-based filtering | ✓ Implemented (frontend tab filter) |
| Auto-creation on events | ✗ Missing (currently manual/stub) |
| Push notifications (FCM/APNs) | ✗ Not planned |

### Admin

| Feature | Status |
|---|---|
| KPI metrics | ✓ Implemented |
| User list + suspend/activate | ✓ Implemented |
| Financial overview + withdrawals | ✓ Implemented |
| Dispute resolution | ✗ Stub (returns 501) |
| Category / city CRUD | ✓ Implemented |
| Audit log | ✗ Missing |
| Bulk notification | ✗ Missing |

### Wallet / Finance

| Feature | Status |
|---|---|
| Per-user wallet auto-creation | ✓ Implemented |
| Top-up | ✓ Implemented |
| Withdraw request | ✓ Implemented |
| Admin approve/reject | ✓ Implemented |
| Transaction history | ✓ Implemented |
| Earnings statistics | ✓ Implemented |
| Payment gateway integration | ✗ Stub (mock top-up) |

---

## 9. Recommendations for Future Phases

### Phase 6 — Trust & Safety

- **Email verification**: Send OTP on register, block login until verified.
- **Dispute resolution**: Implement full dispute ticket flow with admin adjudication. Currently `POST /admin/disputes/:id/resolve` returns 501.
- **Audit logs**: Store every admin action (suspend, approve, resolve) in `audit_logs` table.
- **Password reset**: Email-based reset flow with short-lived token.

### Phase 7 — Real-Time

- **WebSocket chat**: Replace HTTP polling with a WebSocket hub per task. Use Redis pub/sub to fan out across API replicas.
- **Live notifications**: Push badge count updates to topbar without page refresh.
- **Task activity feed**: Stream state changes to all parties (requester, agent) in real-time.

### Phase 8 — Discovery & Matching

- **User skills**: `user_skills` table with structured tags. Filter tasks by skill match.
- **Geolocation**: lat/lng on tasks and user profiles; radius search via PostGIS.
- **Recommendation engine**: Score open tasks by agent skill match + proximity + price range.
- **Task templates**: Predefined task skeletons for common categories.

### Phase 9 — Payments

- **Payment gateway**: Integrate Zarinpal (Iran) or Stripe for real top-up flows.
- **Auto-withdraw**: Scheduled nightly job to process approved withdraw_requests via bank API.
- **Invoice generation**: PDF/email receipt for each escrow release.
- **VAT / tax tracking**: Breakdown of platform fee vs. net agent payment.

### Phase 10 — Growth

- **Reviews on both sides**: Requesters can currently be reviewed too (schema supports it), but the UI only shows reviews received. Add mutual review flow.
- **Referral program**: Referral code field on user, affiliate balance tracking.
- **Agent onboarding checklist**: Guided setup (verify ID, add bank account, pick skills) to increase conversion from registered → active agent.
- **Mobile app**: API is fully REST-based; React Native client would share all existing endpoints.

---

## 10. Technical Debt & Known Issues

| Issue | Severity | Notes |
|---|---|---|
| `receiverId` hardcoded in chat send | Medium | Chat POST requires `receiverId` but client should derive it from task; backend should derive automatically |
| `admin/disputes/:id/resolve` returns 501 | High | Dispute resolution is a core escrow feature |
| No event consumers (RabbitMQ declared but unused) | Medium | All async side-effects (notifications, ledger entries) happen synchronously in HTTP handlers |
| Avatar upload not implemented | Low | `/me/avatar` endpoint missing |
| `wallet/statistics` `totalEarned` field inconsistency | Low | Field name varies between wallet and statistics endpoints |
| Admin chart data is static SVG | Low | Revenue chart in admin dashboard uses hardcoded data; needs real data from `GET /admin/revenue/daily` |
| Chat `receiverId` derivation | Medium | Frontend sends `00000000-…` placeholder; backend should auto-determine receiver from task participants |
