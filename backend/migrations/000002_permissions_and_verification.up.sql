-- Permissions and verification extensions for Phase 3

CREATE TABLE IF NOT EXISTS permissions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role       user_role NOT NULL,
    permission VARCHAR(100) NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
    PRIMARY KEY (role, permission)
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_reset_user ON password_reset_tokens(user_id) WHERE used_at IS NULL;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS verification_status VARCHAR(30) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS verification_reason  TEXT,
    ADD COLUMN IF NOT EXISTS verified_at          TIMESTAMPTZ;

INSERT INTO permissions (name, description) VALUES
    ('task:create', 'Create tasks'),
    ('task:view', 'View tasks'),
    ('task:edit', 'Edit tasks'),
    ('task:delete', 'Delete tasks'),
    ('application:create', 'Submit applications'),
    ('application:accept', 'Accept applications'),
    ('wallet:view', 'View wallet'),
    ('wallet:withdraw', 'Withdraw funds'),
    ('user:manage', 'Manage users'),
    ('admin:dashboard', 'Access admin dashboard'),
    ('review:create', 'Create reviews'),
    ('chat:send', 'Send chat messages'),
    ('notification:view', 'View notifications')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role, permission) VALUES
    ('ADMIN', 'task:create'), ('ADMIN', 'task:view'), ('ADMIN', 'task:edit'), ('ADMIN', 'task:delete'),
    ('ADMIN', 'application:create'), ('ADMIN', 'application:accept'),
    ('ADMIN', 'wallet:view'), ('ADMIN', 'wallet:withdraw'),
    ('ADMIN', 'user:manage'), ('ADMIN', 'admin:dashboard'),
    ('ADMIN', 'review:create'), ('ADMIN', 'chat:send'), ('ADMIN', 'notification:view'),
    ('REQUESTER', 'task:create'), ('REQUESTER', 'task:view'), ('REQUESTER', 'task:edit'), ('REQUESTER', 'task:delete'),
    ('REQUESTER', 'application:accept'), ('REQUESTER', 'wallet:view'), ('REQUESTER', 'wallet:withdraw'),
    ('REQUESTER', 'review:create'), ('REQUESTER', 'chat:send'), ('REQUESTER', 'notification:view'),
    ('AGENT', 'task:view'), ('AGENT', 'application:create'),
    ('AGENT', 'wallet:view'), ('AGENT', 'wallet:withdraw'),
    ('AGENT', 'review:create'), ('AGENT', 'chat:send'), ('AGENT', 'notification:view')
ON CONFLICT DO NOTHING;
