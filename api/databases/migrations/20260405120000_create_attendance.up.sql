CREATE TABLE IF NOT EXISTS office_configs (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(200) NOT NULL DEFAULT 'Main Office',
    latitude        DOUBLE PRECISION NOT NULL DEFAULT 0,
    longitude       DOUBLE PRECISION NOT NULL DEFAULT 0,
    radius_meters   DOUBLE PRECISION NOT NULL DEFAULT 100,
    allowed_ips     TEXT[] NOT NULL DEFAULT '{}',
    work_start_time VARCHAR(8) NOT NULL DEFAULT '09:00:00',
    work_end_time   VARCHAR(8) NOT NULL DEFAULT '18:00:00',
    work_days       JSONB NOT NULL DEFAULT '[1,2,3,4,5]'::jsonb,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_office_configs_active ON office_configs (is_active) WHERE is_active = true;

CREATE TABLE IF NOT EXISTS attendance_records (
    id                 BIGSERIAL PRIMARY KEY,
    user_id            BIGINT NOT NULL,
    office_config_id   BIGINT NOT NULL,
    attendance_date    DATE NOT NULL,
    check_in_at        TIMESTAMPTZ,
    check_out_at       TIMESTAMPTZ,
    check_in_lat       DOUBLE PRECISION,
    check_in_lng       DOUBLE PRECISION,
    check_in_method    VARCHAR(20) NOT NULL DEFAULT '',
    check_in_ip        VARCHAR(64) NOT NULL DEFAULT '',
    is_late            BOOLEAN NOT NULL DEFAULT false,
    early_checkout     BOOLEAN NOT NULL DEFAULT false,
    status             VARCHAR(20) NOT NULL DEFAULT 'absent',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_attendance_user_date UNIQUE (user_id, attendance_date),
    CONSTRAINT fk_attendance_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_attendance_office FOREIGN KEY (office_config_id) REFERENCES office_configs(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_attendance_records_user_id ON attendance_records (user_id);
CREATE INDEX IF NOT EXISTS idx_attendance_records_date ON attendance_records (attendance_date);
CREATE INDEX IF NOT EXISTS idx_attendance_records_user_id_id ON attendance_records (user_id, id);
