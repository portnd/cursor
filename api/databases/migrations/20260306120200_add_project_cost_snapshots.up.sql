-- Versioned cost estimates per project (immutable once APPROVED)
CREATE TABLE IF NOT EXISTS project_cost_snapshots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 1,
    total_labor_cost DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_expenses DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_cost DECIMAL(15,2) NOT NULL DEFAULT 0,
    suggested_price DECIMAL(15,2) NOT NULL DEFAULT 0,
    profit_margin DECIMAL(5,2) NOT NULL DEFAULT 0.25,
    risk_buffer DECIMAL(5,2) NOT NULL DEFAULT 0.10,
    estimated_hours DECIMAL(10,2) NOT NULL DEFAULT 0,
    estimated_days DECIMAL(10,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    breakdown JSONB DEFAULT '[]',
    notes TEXT,
    valid_until DATE,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_snapshot_status CHECK (status IN ('DRAFT', 'APPROVED', 'SENT_TO_CLIENT'))
);

CREATE INDEX IF NOT EXISTS idx_project_cost_snapshots_project_id ON project_cost_snapshots(project_id);
CREATE INDEX IF NOT EXISTS idx_project_cost_snapshots_status ON project_cost_snapshots(status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_project_cost_snapshots_project_version ON project_cost_snapshots(project_id, version);
