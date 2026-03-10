-- The Handover: simplify submissions table — remove AI fields, add reference_url + note
-- Drop constraints first, then columns

ALTER TABLE submissions DROP CONSTRAINT IF EXISTS check_ai_verdict;
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS check_ai_score;
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS check_commit_hash;

DROP INDEX IF EXISTS idx_submissions_task_commit;

ALTER TABLE submissions
    DROP COLUMN IF EXISTS commit_hash,
    DROP COLUMN IF EXISTS diff,
    DROP COLUMN IF EXISTS ai_verdict,
    DROP COLUMN IF EXISTS ai_score,
    DROP COLUMN IF EXISTS ai_feedback,
    DROP COLUMN IF EXISTS is_overridden;

ALTER TABLE submissions
    ADD COLUMN IF NOT EXISTS reference_url VARCHAR(512),
    ADD COLUMN IF NOT EXISTS note TEXT;
