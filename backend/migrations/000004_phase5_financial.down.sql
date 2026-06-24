-- Phase 5 rollback.

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS withdraw_requests;
DROP TYPE IF EXISTS withdraw_status;
DROP TABLE IF EXISTS revenue_records;
DROP TABLE IF EXISTS ledger_entries;

ALTER TABLE transactions DROP COLUMN IF EXISTS currency;

ALTER TABLE wallets DROP COLUMN IF EXISTS currency;
ALTER TABLE wallets DROP COLUMN IF EXISTS total_spent;
ALTER TABLE wallets DROP COLUMN IF EXISTS total_earned;
ALTER TABLE wallets DROP COLUMN IF EXISTS pending_withdraw_balance;
