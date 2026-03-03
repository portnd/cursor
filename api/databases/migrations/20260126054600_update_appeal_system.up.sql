-- Update Appeal System to use submission_id instead of task_id
-- Ensure UUID extension exists (this migration may run before init_sentinel_schema)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop existing appeals table
DROP TABLE IF EXISTS appeals;

-- Recreate appeals table with new schema
CREATE TABLE appeals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    developer_id INT NOT NULL REFERENCES users(id),
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
    resolver_id INT REFERENCES users(id),
    resolver_note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add index on submission_id for faster lookups
CREATE INDEX idx_appeals_submission_id ON appeals(submission_id);

-- Add is_overridden column to submissions table
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS is_overridden BOOLEAN DEFAULT FALSE;

-- Add index on is_overridden for queries
CREATE INDEX IF NOT EXISTS idx_submissions_is_overridden ON submissions(is_overridden);
