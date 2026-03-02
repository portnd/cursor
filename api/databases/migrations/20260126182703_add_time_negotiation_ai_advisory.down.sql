-- Remove AI Advisory fields for Time Negotiation

ALTER TABLE tasks DROP COLUMN IF EXISTS negotiation_ai_recommendation;
ALTER TABLE tasks DROP COLUMN IF EXISTS negotiation_ai_confidence;
ALTER TABLE tasks DROP COLUMN IF EXISTS negotiation_ai_reasoning;
