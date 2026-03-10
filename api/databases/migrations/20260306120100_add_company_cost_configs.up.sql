-- Singleton: company-wide cost calculation settings
CREATE TABLE IF NOT EXISTS company_cost_configs (
    id BIGSERIAL PRIMARY KEY,
    working_days_per_month INT NOT NULL DEFAULT 22,
    working_hours_per_day INT NOT NULL DEFAULT 8,
    overhead_multiplier DECIMAL(5,2) NOT NULL DEFAULT 1.30,
    default_profit_margin DECIMAL(5,2) NOT NULL DEFAULT 0.25,
    default_risk_buffer DECIMAL(5,2) NOT NULL DEFAULT 0.10,
    currency VARCHAR(3) NOT NULL DEFAULT 'THB',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Ensure single row: insert default if empty
INSERT INTO company_cost_configs (id, working_days_per_month, working_hours_per_day, overhead_multiplier, default_profit_margin, default_risk_buffer, currency)
SELECT 1, 22, 8, 1.30, 0.25, 0.10, 'THB'
WHERE NOT EXISTS (SELECT 1 FROM company_cost_configs LIMIT 1);
