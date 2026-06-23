ALTER TABLE users
    DROP COLUMN IF EXISTS verified_at,
    DROP COLUMN IF EXISTS verification_reason,
    DROP COLUMN IF EXISTS verification_status;

DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
