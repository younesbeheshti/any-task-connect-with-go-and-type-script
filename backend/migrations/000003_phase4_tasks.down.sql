-- Rollback Phase 4 migration
ALTER TABLE tasks DROP COLUMN IF EXISTS view_count;
ALTER TABLE tasks DROP COLUMN IF EXISTS application_count;
ALTER TABLE applications DROP COLUMN IF EXISTS eta;
