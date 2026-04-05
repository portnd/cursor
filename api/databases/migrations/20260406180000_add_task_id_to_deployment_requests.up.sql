-- Link deployment requests to sentinel tasks so that marking a request as DEPLOYED
-- automatically advances the task from WAIT_FOR_DEPLOY → READY_FOR_UAT.
ALTER TABLE deployment_requests ADD COLUMN IF NOT EXISTS task_id UUID REFERENCES tasks(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_deployment_requests_task_id ON deployment_requests(task_id);
