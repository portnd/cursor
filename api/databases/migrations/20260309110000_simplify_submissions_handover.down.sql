-- Rollback: restore AI fields to submissions table

ALTER TABLE submissions
    DROP COLUMN IF EXISTS reference_url,
    DROP COLUMN IF EXISTS note;

ALTER TABLE submissions
    ADD COLUMN IF NOT EXISTS commit_hash VARCHAR(64),
    ADD COLUMN IF NOT EXISTS diff TEXT,
    ADD COLUMN IF NOT EXISTS ai_verdict VARCHAR(20),
    ADD COLUMN IF NOT EXISTS ai_score INT,
    ADD COLUMN IF NOT EXISTS ai_feedback JSONB DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS is_overridden BOOLEAN DEFAULT FALSE;

ALTER TABLE submissions
    ADD CONSTRAINT check_ai_verdict CHECK (ai_verdict IN ('PASS', 'FAIL', 'PENDING')),
    ADD CONSTRAINT check_ai_score CHECK (ai_score IS NULL OR (ai_score >= 0 AND ai_score <= 100)),
    ADD CONSTRAINT check_commit_hash CHECK (commit_hash IS NULL OR (LENGTH(commit_hash) >= 7 AND LENGTH(commit_hash) <= 64));

CREATE UNIQUE INDEX IF NOT EXISTS idx_submissions_task_commit ON submissions(task_id, commit_hash)
    WHERE commit_hash IS NOT NULL;
