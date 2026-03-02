-- Add Time Negotiation fields to tasks table

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS negotiation_status VARCHAR(20) DEFAULT 'NONE' CHECK (negotiation_status IN ('NONE', 'PENDING', 'APPROVED', 'REJECTED'));
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS proposed_minutes INT DEFAULT 0;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS negotiation_reason TEXT;

-- Add index on negotiation_status for filtering
CREATE INDEX IF NOT EXISTS idx_tasks_negotiation_status ON tasks(negotiation_status);

-- Add comments for documentation
COMMENT ON COLUMN tasks.negotiation_status IS 'Time negotiation status: NONE (no negotiation), PENDING (awaiting review), APPROVED (accepted), REJECTED (denied)';
COMMENT ON COLUMN tasks.proposed_minutes IS 'Developer proposed time in minutes (if they dispute AI estimate)';
COMMENT ON COLUMN tasks.negotiation_reason IS 'Developer explanation for why they need more time than AI estimated';
