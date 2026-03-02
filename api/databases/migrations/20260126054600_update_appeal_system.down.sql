-- Rollback Appeal System changes

-- Drop new indexes
DROP INDEX IF EXISTS idx_appeals_submission_id;
DROP INDEX IF EXISTS idx_submissions_is_overridden;

-- Remove is_overridden column from submissions
ALTER TABLE submissions DROP COLUMN IF EXISTS is_overridden;

-- Drop new appeals table
DROP TABLE IF EXISTS appeals;

-- Recreate old appeals table structure (if needed)
CREATE TABLE appeals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id),
    dev_id INT NOT NULL REFERENCES users(id),
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING',
    reviewed_by INT REFERENCES users(id),
    admin_comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
