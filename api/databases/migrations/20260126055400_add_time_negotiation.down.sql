-- Rollback Time Negotiation fields

DROP INDEX IF EXISTS idx_tasks_negotiation_status;

ALTER TABLE tasks DROP COLUMN IF EXISTS negotiation_reason;
ALTER TABLE tasks DROP COLUMN IF EXISTS proposed_minutes;
ALTER TABLE tasks DROP COLUMN IF EXISTS negotiation_status;
