-- Phase 4: add REJECTED status to task_status enum, add application ETA field
-- The core tables already exist from migration 000001.
-- This migration adds the REJECTED status and view_count column.

ALTER TYPE task_status ADD VALUE IF NOT EXISTS 'REJECTED';

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS view_count INT NOT NULL DEFAULT 0;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS application_count INT NOT NULL DEFAULT 0;

-- Add ETA column to applications if not present
ALTER TABLE applications ADD COLUMN IF NOT EXISTS eta VARCHAR(50) NOT NULL DEFAULT '';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS proposed_price BIGINT;
