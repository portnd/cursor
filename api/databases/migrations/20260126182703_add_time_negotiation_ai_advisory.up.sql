-- Add AI Advisory fields for Time Negotiation

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS negotiation_ai_recommendation VARCHAR(20) CHECK (negotiation_ai_recommendation IN ('APPROVE', 'REJECT', ''));
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS negotiation_ai_confidence INT DEFAULT 0 CHECK (negotiation_ai_confidence >= 0 AND negotiation_ai_confidence <= 100);
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS negotiation_ai_reasoning TEXT;

-- Add comments
COMMENT ON COLUMN tasks.negotiation_ai_recommendation IS 'AI recommendation for time negotiation: APPROVE (developer estimate is reasonable), REJECT (keep AI estimate)';
COMMENT ON COLUMN tasks.negotiation_ai_confidence IS 'AI confidence level (0-100) in the recommendation';
COMMENT ON COLUMN tasks.negotiation_ai_reasoning IS 'AI explanation for why developer request should be approved or rejected';
