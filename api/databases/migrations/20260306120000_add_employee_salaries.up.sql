-- Employee salaries: effective-dated monthly salary per user for cost estimation
CREATE TABLE IF NOT EXISTS employee_salaries (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    monthly_salary DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'THB',
    effective_from DATE NOT NULL,
    effective_to DATE,
    employment_type VARCHAR(20) NOT NULL DEFAULT 'FULLTIME',
    cost_per_minute DECIMAL(10,6),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_employment_type CHECK (employment_type IN ('FULLTIME', 'PARTTIME', 'CONTRACTOR'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_employee_salaries_user_effective ON employee_salaries(user_id, effective_from);
CREATE INDEX IF NOT EXISTS idx_employee_salaries_user_id ON employee_salaries(user_id);
CREATE INDEX IF NOT EXISTS idx_employee_salaries_effective ON employee_salaries(effective_from, effective_to);
