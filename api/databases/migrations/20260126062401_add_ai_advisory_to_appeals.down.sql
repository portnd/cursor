-- Rollback: Remove AI Advisory System from Appeals

-- Drop indexes
DROP INDEX IF EXISTS idx_appeals_ai_confidence;
DROP INDEX IF EXISTS idx_appeals_ai_recommendation;

-- Remove diff from submissions
ALTER TABLE submissions
DROP COLUMN IF EXISTS diff;

-- Remove AI Advisory fields from appeals
ALTER TABLE appeals
DROP COLUMN IF EXISTS ai_reasoning,
DROP COLUMN IF EXISTS ai_confidence,
DROP COLUMN IF EXISTS ai_recommendation;
