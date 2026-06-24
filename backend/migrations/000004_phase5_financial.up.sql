-- Phase 5: Financial domain — extend wallets/transactions, add ledger, revenue, withdraw, audit tables.

ALTER TABLE wallets ADD COLUMN IF NOT EXISTS pending_withdraw_balance BIGINT NOT NULL DEFAULT 0 CHECK (pending_withdraw_balance >= 0);
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS total_earned BIGINT NOT NULL DEFAULT 0;
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS total_spent BIGINT NOT NULL DEFAULT 0;
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS currency VARCHAR(10) NOT NULL DEFAULT 'IRT';

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS currency VARCHAR(10) NOT NULL DEFAULT 'IRT';

ALTER TYPE transaction_type ADD VALUE IF NOT EXISTS 'PAYMENT';
ALTER TYPE transaction_type ADD VALUE IF NOT EXISTS 'COMMISSION';
ALTER TYPE transaction_type ADD VALUE IF NOT EXISTS 'ADJUSTMENT';
ALTER TYPE transaction_type ADD VALUE IF NOT EXISTS 'DEPOSIT';

CREATE TABLE IF NOT EXISTS ledger_entries (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    account        VARCHAR(50) NOT NULL,
    debit          BIGINT NOT NULL DEFAULT 0 CHECK (debit >= 0),
    credit         BIGINT NOT NULL DEFAULT 0 CHECK (credit >= 0),
    balance_after  BIGINT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ledger_transaction ON ledger_entries(transaction_id);
CREATE INDEX IF NOT EXISTS idx_ledger_account ON ledger_entries(account, created_at DESC);

CREATE TABLE IF NOT EXISTS revenue_records (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id    UUID REFERENCES tasks(id),
    source     VARCHAR(50) NOT NULL DEFAULT 'commission',
    amount     BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_revenue_created ON revenue_records(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_revenue_task ON revenue_records(task_id) WHERE task_id IS NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'withdraw_status') THEN
        CREATE TYPE withdraw_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED', 'PAID');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS withdraw_requests (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(id),
    wallet_id      UUID NOT NULL REFERENCES wallets(id),
    amount         BIGINT NOT NULL CHECK (amount > 0),
    status         withdraw_status NOT NULL DEFAULT 'PENDING',
    bank_account   VARCHAR(30),
    sheba_number   VARCHAR(26),
    description    TEXT NOT NULL DEFAULT '',
    approved_by    UUID REFERENCES users(id),
    approved_at    TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_withdraw_user ON withdraw_requests(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_withdraw_status ON withdraw_requests(status, created_at DESC);

CREATE TABLE IF NOT EXISTS audit_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type  VARCHAR(50) NOT NULL,
    entity_id    UUID,
    action       VARCHAR(50) NOT NULL,
    performed_by UUID REFERENCES users(id),
    old_value    JSONB,
    new_value    JSONB,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_logs(entity_type, entity_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_performer ON audit_logs(performed_by, created_at DESC);
