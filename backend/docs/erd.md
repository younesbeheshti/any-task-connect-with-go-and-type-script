# TaskBridge — Entity Relationship Diagram

## ER Diagram (Mermaid)

```mermaid
erDiagram
    USERS ||--o| WALLETS : has
    USERS ||--o{ TASKS : "creates (requester)"
    USERS ||--o{ TASKS : "assigned (agent)"
    USERS ||--o{ APPLICATIONS : submits
    USERS ||--o{ NOTIFICATIONS : receives
    USERS ||--o{ REVIEWS : "writes (reviewer)"
    USERS ||--o{ REVIEWS : "receives (reviewed)"
    USERS ||--o{ CHAT_MESSAGES : sends

    CATEGORIES ||--o{ TASKS : categorizes
    CITIES ||--o{ TASKS : locates

    TASKS ||--o{ APPLICATIONS : receives
    TASKS ||--o{ TRANSACTIONS : references
    TASKS ||--o{ CHAT_MESSAGES : context
    TASKS ||--o{ REVIEWS : rated_via
    TASKS ||--o{ TASK_TIMELINE : tracks

    WALLETS ||--o{ TRANSACTIONS : records

    USERS {
        uuid id PK
        varchar full_name
        varchar phone UK
        varchar email UK
        varchar password_hash
        enum role
        varchar national_id
        varchar avatar
        decimal rating
        int rating_count
        int completed_tasks
        boolean is_verified
        boolean is_active
        bigint wallet_balance
        bigint locked_balance
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    WALLETS {
        uuid id PK
        uuid user_id FK UK
        bigint available_balance
        bigint locked_balance
        timestamptz created_at
        timestamptz updated_at
    }

    TASKS {
        uuid id PK
        varchar public_id UK
        varchar title
        text description
        uuid category_id FK
        uuid city_id FK
        bigint budget
        enum status
        date deadline
        uuid requester_id FK
        uuid assigned_agent_id FK
        jsonb attachment_urls
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    APPLICATIONS {
        uuid id PK
        uuid task_id FK
        uuid agent_id FK
        text proposal_message
        timestamptz expected_completion_time
        bigint proposed_price
        enum status
        timestamptz created_at
        timestamptz updated_at
    }

    TRANSACTIONS {
        uuid id PK
        uuid wallet_id FK
        uuid task_id FK
        bigint amount
        enum type
        enum status
        varchar reference_number UK
        timestamptz created_at
    }

    CHAT_MESSAGES {
        uuid id PK
        uuid task_id FK
        uuid sender_id FK
        uuid receiver_id FK
        text message
        jsonb attachment
        boolean seen
        timestamptz created_at
    }

    NOTIFICATIONS {
        uuid id PK
        uuid user_id FK
        varchar title
        text body
        enum type
        boolean is_read
        varchar deep_link
        timestamptz created_at
    }

    REVIEWS {
        uuid id PK
        uuid task_id FK
        uuid reviewer_id FK
        uuid reviewed_user_id FK
        smallint rating
        text comment
        timestamptz created_at
    }

    CATEGORIES {
        uuid id PK
        varchar title UK
        varchar icon
        timestamptz created_at
    }

    CITIES {
        uuid id PK
        varchar title
        varchar province
        timestamptz created_at
    }

    TASK_TIMELINE {
        uuid id PK
        uuid task_id FK
        enum from_status
        enum to_status
        uuid actor_id FK
        text note
        timestamptz created_at
    }
```

## Supporting Tables (Phase 2+)

| Table | Purpose |
|-------|---------|
| `refresh_tokens` | Hashed refresh token storage |
| `password_reset_tokens` | Forgot password flow |
| `email_verifications` | Email OTP / link verification |
| `phone_verifications` | Phone OTP |
| `files` | Uploaded attachment metadata |
| `disputes` | Admin dispute resolution |
| `bank_cards` | Wallet top-up cards |

## Key Constraints

- **Users**: unique `phone`, optional unique `email`
- **Tasks**: unique `public_id` (e.g. `TB-1042`) for API display
- **Applications**: unique `(task_id, agent_id)` — one application per agent per task
- **Reviews**: unique `(task_id, reviewer_id)` — one review per participant per task
- **Wallets**: one wallet per user (`user_id` unique)
- **Transactions**: immutable ledger; status changes create new records or append audit

## Index Strategy

| Table | Index | Reason |
|-------|-------|--------|
| `tasks` | `(status, city_id, category_id)` | Agent task browsing |
| `tasks` | `(requester_id)` | Requester dashboard |
| `tasks` | `(assigned_agent_id)` | Agent assigned tasks |
| `applications` | `(task_id, status)` | Application listing |
| `notifications` | `(user_id, is_read, created_at DESC)` | Unread feed |
| `chat_messages` | `(task_id, created_at)` | Message pagination |
| `transactions` | `(wallet_id, created_at DESC)` | Wallet history |
