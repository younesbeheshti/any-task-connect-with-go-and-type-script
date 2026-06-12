-- TaskBridge initial schema
-- Currency: Toman (bigint, no decimals)

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enums
CREATE TYPE user_role AS ENUM ('ADMIN', 'REQUESTER', 'AGENT');
CREATE TYPE task_status AS ENUM (
    'CREATED', 'OPEN', 'ASSIGNED', 'IN_PROGRESS', 'COMPLETED',
    'WAITING_FOR_VERIFICATION', 'VERIFIED', 'PAID', 'CANCELLED'
);
CREATE TYPE application_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED', 'WITHDRAWN');
CREATE TYPE transaction_type AS ENUM (
    'ESCROW_LOCK', 'ESCROW_RELEASE', 'TOPUP', 'REFUND', 'WITHDRAW', 'FEE', 'TRANSFER'
);
CREATE TYPE transaction_status AS ENUM ('PENDING', 'COMPLETED', 'FAILED', 'LOCKED', 'RELEASED');
CREATE TYPE notification_type AS ENUM (
    'TASK_ASSIGNED', 'TASK_COMPLETED', 'PAYMENT', 'APPLICATION',
    'MESSAGE', 'DISPUTE', 'SYSTEM', 'REVIEW'
);

-- Categories
CREATE TABLE categories (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title      VARCHAR(100) NOT NULL UNIQUE,
    icon       VARCHAR(50)  NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Cities
CREATE TABLE cities (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title      VARCHAR(100) NOT NULL,
    province   VARCHAR(100) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (title, province)
);

-- Users
CREATE TABLE users (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name        VARCHAR(200)  NOT NULL,
    phone            VARCHAR(20)   NOT NULL UNIQUE,
    email            VARCHAR(255)  UNIQUE,
    password_hash    VARCHAR(255)  NOT NULL,
    role             user_role     NOT NULL DEFAULT 'REQUESTER',
    national_id      VARCHAR(20),
    avatar           TEXT,
    bio              TEXT,
    city_id          UUID REFERENCES cities(id),
    rating           NUMERIC(3,2)  NOT NULL DEFAULT 0 CHECK (rating >= 0 AND rating <= 5),
    rating_count     INT           NOT NULL DEFAULT 0,
    completed_tasks  INT           NOT NULL DEFAULT 0,
    is_verified      BOOLEAN       NOT NULL DEFAULT FALSE,
    is_active        BOOLEAN       NOT NULL DEFAULT TRUE,
    phone_verified   BOOLEAN       NOT NULL DEFAULT FALSE,
    email_verified   BOOLEAN       NOT NULL DEFAULT FALSE,
    national_id_verified BOOLEAN   NOT NULL DEFAULT FALSE,
    verification_level VARCHAR(20) NOT NULL DEFAULT 'none',
    wallet_balance   BIGINT        NOT NULL DEFAULT 0 CHECK (wallet_balance >= 0),
    locked_balance   BIGINT        NOT NULL DEFAULT 0 CHECK (locked_balance >= 0),
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_city ON users(city_id) WHERE deleted_at IS NULL;

-- Wallets (normalized; synced with user denormalized balances via triggers or service layer)
CREATE TABLE wallets (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE RESTRICT,
    available_balance BIGINT NOT NULL DEFAULT 0 CHECK (available_balance >= 0),
    locked_balance    BIGINT NOT NULL DEFAULT 0 CHECK (locked_balance >= 0),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tasks
CREATE TABLE tasks (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id          VARCHAR(20) NOT NULL UNIQUE,
    title              VARCHAR(120) NOT NULL,
    description        TEXT NOT NULL,
    category_id        UUID NOT NULL REFERENCES categories(id),
    city_id            UUID NOT NULL REFERENCES cities(id),
    budget             BIGINT NOT NULL CHECK (budget > 0),
    escrow_fee         BIGINT NOT NULL DEFAULT 0 CHECK (escrow_fee >= 0),
    status             task_status NOT NULL DEFAULT 'CREATED',
    deadline           DATE NOT NULL,
    requester_id       UUID NOT NULL REFERENCES users(id),
    assigned_agent_id  UUID REFERENCES users(id),
    attachment_urls    JSONB NOT NULL DEFAULT '[]',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);

CREATE INDEX idx_tasks_status_city_category ON tasks(status, city_id, category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tasks_requester ON tasks(requester_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tasks_agent ON tasks(assigned_agent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tasks_created ON tasks(created_at DESC) WHERE deleted_at IS NULL;

-- Task timeline (audit)
CREATE TABLE task_timeline (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id     UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    from_status task_status,
    to_status   task_status NOT NULL,
    actor_id    UUID REFERENCES users(id),
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_task_timeline_task ON task_timeline(task_id, created_at);

-- Applications
CREATE TABLE applications (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id                  UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    agent_id                 UUID NOT NULL REFERENCES users(id),
    proposal_message         TEXT NOT NULL DEFAULT '',
    expected_completion_time TIMESTAMPTZ,
    proposed_price           BIGINT CHECK (proposed_price IS NULL OR proposed_price > 0),
    status                   application_status NOT NULL DEFAULT 'PENDING',
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (task_id, agent_id)
);

CREATE INDEX idx_applications_task_status ON applications(task_id, status);
CREATE INDEX idx_applications_agent ON applications(agent_id);

-- Transactions (immutable ledger)
CREATE TABLE transactions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id        UUID NOT NULL REFERENCES wallets(id),
    task_id          UUID REFERENCES tasks(id),
    amount           BIGINT NOT NULL,
    type             transaction_type NOT NULL,
    status           transaction_status NOT NULL DEFAULT 'PENDING',
    reference_number VARCHAR(50) NOT NULL UNIQUE,
    description      TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_wallet ON transactions(wallet_id, created_at DESC);
CREATE INDEX idx_transactions_task ON transactions(task_id) WHERE task_id IS NOT NULL;

-- Chat messages
CREATE TABLE chat_messages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id     UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    sender_id   UUID NOT NULL REFERENCES users(id),
    receiver_id UUID NOT NULL REFERENCES users(id),
    message     TEXT NOT NULL DEFAULT '',
    attachment  JSONB,
    seen        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chat_messages_task ON chat_messages(task_id, created_at DESC);
CREATE INDEX idx_chat_messages_receiver_unread ON chat_messages(receiver_id, seen) WHERE seen = FALSE;

-- Notifications
CREATE TABLE notifications (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    body       TEXT NOT NULL DEFAULT '',
    type       notification_type NOT NULL,
    is_read    BOOLEAN NOT NULL DEFAULT FALSE,
    deep_link  VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications(user_id, is_read, created_at DESC);

-- Reviews
CREATE TABLE reviews (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id          UUID NOT NULL REFERENCES tasks(id),
    reviewer_id      UUID NOT NULL REFERENCES users(id),
    reviewed_user_id UUID NOT NULL REFERENCES users(id),
    rating           SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment          TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (task_id, reviewer_id)
);

CREATE INDEX idx_reviews_reviewed_user ON reviews(reviewed_user_id);

-- Refresh tokens
CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id) WHERE revoked_at IS NULL;

-- Seed categories (Persian labels from frontend)
INSERT INTO categories (title, icon) VALUES
    ('اداری', 'FileText'),
    ('پزشکی', 'Stethoscope'),
    ('دانشگاهی', 'GraduationCap'),
    ('حقوقی', 'Scale'),
    ('دولتی', 'Landmark'),
    ('پیک و تحویل', 'Package'),
    ('خرید', 'ShoppingBag'),
    ('سایر', 'Sparkles');

-- Seed cities
INSERT INTO cities (title, province) VALUES
    ('تهران', 'تهران'),
    ('مشهد', 'خراسان رضوی'),
    ('اصفهان', 'اصفهان'),
    ('شیراز', 'فارس'),
    ('تبریز', 'آذربایجان شرقی'),
    ('کرج', 'البرز'),
    ('اهواز', 'خوزستان'),
    ('قم', 'قم'),
    ('رشت', 'گیلان'),
    ('کرمان', 'کرمان');
