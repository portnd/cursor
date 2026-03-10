-- Revert: rename estimated_minutes back to ai_estimated_minutes
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS check_estimated_minutes;
ALTER TABLE tasks RENAME COLUMN estimated_minutes TO ai_estimated_minutes;
ALTER TABLE tasks ADD CONSTRAINT check_ai_estimated_minutes CHECK (ai_estimated_minutes IS NULL OR ai_estimated_minutes >= 0);
