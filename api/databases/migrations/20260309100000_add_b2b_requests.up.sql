-- Internal B2B Outsource Requests table
-- Tracks cross-team work requests: Team A asks Team B to handle a task

CREATE TABLE IF NOT EXISTS b2b_requests (
  id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  title               VARCHAR(255) NOT NULL,
  description         TEXT,
  estimated_minutes   INT         NOT NULL DEFAULT 0,
  proposed_minutes    INT         NOT NULL DEFAULT 0,
  negotiation_reason  TEXT,
  status              VARCHAR(30) NOT NULL DEFAULT 'PENDING'
                        CHECK (status IN ('PENDING', 'COUNTER_OFFERED', 'ACCEPTED', 'REJECTED')),

  requester_team_id   INT         NOT NULL REFERENCES teams(id),
  target_team_id      INT         NOT NULL REFERENCES teams(id),
  requester_user_id   INT         NOT NULL,

  -- Populated after status becomes ACCEPTED (the task created in target team's project)
  created_task_id     UUID        REFERENCES tasks(id) ON DELETE SET NULL,

  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_b2b_requests_target_team    ON b2b_requests(target_team_id);
CREATE INDEX IF NOT EXISTS idx_b2b_requests_requester_team ON b2b_requests(requester_team_id);
CREATE INDEX IF NOT EXISTS idx_b2b_requests_status         ON b2b_requests(status);
