-- Migrate data from ai_estimated_minutes → estimated_minutes and drop the old column
-- Note: estimated_minutes column already exists (added by GORM AutoMigrate when entity was updated)

-- Copy any non-zero values from old column to new column (where new column is still 0/null)
UPDATE tasks
SET estimated_minutes = ai_estimated_minutes
WHERE (estimated_minutes IS NULL OR estimated_minutes = 0)
  AND ai_estimated_minutes IS NOT NULL
  AND ai_estimated_minutes > 0;

-- Drop the old column
ALTER TABLE tasks DROP COLUMN IF EXISTS ai_estimated_minutes;

-- Drop old constraint (may not exist if already dropped)
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS check_ai_estimated_minutes;

-- Ensure new constraint exists
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS check_estimated_minutes;
ALTER TABLE tasks ADD CONSTRAINT check_estimated_minutes CHECK (estimated_minutes IS NULL OR estimated_minutes >= 0);
