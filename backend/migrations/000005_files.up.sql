-- File uploads: store metadata for files written to local disk.

CREATE TABLE IF NOT EXISTS files (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_name VARCHAR(512) NOT NULL,
    stored_name   VARCHAR(512) NOT NULL UNIQUE,
    path          VARCHAR(1024) NOT NULL,
    size          BIGINT NOT NULL CHECK (size >= 0),
    mime_type     VARCHAR(255) NOT NULL,
    uploaded_by   UUID NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_files_uploaded_by ON files(uploaded_by, created_at DESC);
