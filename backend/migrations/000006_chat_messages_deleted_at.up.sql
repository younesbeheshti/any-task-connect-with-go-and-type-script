-- chat_messages uses GORM soft-delete (DeletedAt) and read receipts (ReadAt), and
-- the chat queries filter on deleted_at, but both columns were missing from the
-- initial schema. Add them.

ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS read_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_chat_messages_deleted_at ON chat_messages(deleted_at);
