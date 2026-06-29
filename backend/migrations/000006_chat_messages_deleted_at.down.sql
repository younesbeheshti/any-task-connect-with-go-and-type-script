DROP INDEX IF EXISTS idx_chat_messages_deleted_at;
ALTER TABLE chat_messages DROP COLUMN IF EXISTS read_at;
ALTER TABLE chat_messages DROP COLUMN IF EXISTS deleted_at;
