-- Migration: Add AI Advisory System to Appeals
-- This adds AI recommendation fields to help CEO/PM make decisions

-- Add AI Advisory fields to appeals table
ALTER TABLE appeals
ADD COLUMN ai_recommendation TEXT,
ADD COLUMN ai_confidence INTEGER DEFAULT 0,
ADD COLUMN ai_reasoning TEXT;

-- Add diff storage to submissions for appeal analysis
ALTER TABLE submissions
ADD COLUMN diff TEXT;

-- Add comments for documentation
COMMENT ON COLUMN appeals.ai_recommendation IS 'AI recommendation: OVERTURN (approve) or UPHOLD (reject)';
COMMENT ON COLUMN appeals.ai_confidence IS 'AI confidence score (0-100)';
COMMENT ON COLUMN appeals.ai_reasoning IS 'AI explanation for CEO/PM to consider';
COMMENT ON COLUMN submissions.diff IS 'Code diff for appeal analysis';

-- Create index for faster appeal analysis queries
CREATE INDEX IF NOT EXISTS idx_appeals_ai_recommendation ON appeals(ai_recommendation);
CREATE INDEX IF NOT EXISTS idx_appeals_ai_confidence ON appeals(ai_confidence);
