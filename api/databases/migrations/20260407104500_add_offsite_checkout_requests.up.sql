CREATE TABLE IF NOT EXISTS offsite_checkout_requests (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    office_config_id BIGINT NOT NULL,
    attendance_date DATE NOT NULL,
    request_lat     DOUBLE PRECISION NOT NULL,
    request_lng     DOUBLE PRECISION NOT NULL,
    reason          TEXT NOT NULL DEFAULT '',
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    approver_id     BIGINT,
    approver_note   TEXT NOT NULL DEFAULT '',
    requested_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_offsite_checkout_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_offsite_checkout_office FOREIGN KEY (office_config_id) REFERENCES office_configs(id) ON DELETE RESTRICT,
    CONSTRAINT fk_offsite_checkout_approver FOREIGN KEY (approver_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_offsite_checkout_user_date ON offsite_checkout_requests (user_id, attendance_date);
CREATE INDEX IF NOT EXISTS idx_offsite_checkout_status ON offsite_checkout_requests (status);
