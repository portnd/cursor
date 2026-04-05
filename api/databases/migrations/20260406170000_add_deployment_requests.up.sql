-- Deployment requests: engineers submit PR/branch for Chief Engineer review & merge/deploy
CREATE TABLE IF NOT EXISTS deployment_requests (
    id            BIGSERIAL PRIMARY KEY,
    title         VARCHAR(200)             NOT NULL,
    description   TEXT,
    branch        VARCHAR(300)             NOT NULL,
    pr_url        VARCHAR(512),
    environment   VARCHAR(20)              NOT NULL DEFAULT 'STAGING',    -- STAGING | PRODUCTION
    status        VARCHAR(20)              NOT NULL DEFAULT 'PENDING',    -- PENDING | REVIEWING | APPROVED | REJECTED | DEPLOYED
    requester_id  BIGINT                   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reviewer_id   BIGINT                   REFERENCES users(id) ON DELETE SET NULL,
    task_ref      VARCHAR(200),                                           -- optional task title / ID reference
    rejection_reason TEXT,
    review_notes  TEXT,
    deployed_at   TIMESTAMP WITH TIME ZONE,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_deployment_requests_status      ON deployment_requests(status);
CREATE INDEX IF NOT EXISTS idx_deployment_requests_requester   ON deployment_requests(requester_id);
CREATE INDEX IF NOT EXISTS idx_deployment_requests_reviewer    ON deployment_requests(reviewer_id);
CREATE INDEX IF NOT EXISTS idx_deployment_requests_created     ON deployment_requests(created_at DESC);
