-- Key-value app settings (e.g. feature flags)
CREATE TABLE IF NOT EXISTS app_settings (
    key VARCHAR(64) PRIMARY KEY,
    value TEXT NOT NULL
);

-- Default: teams feature enabled
INSERT INTO app_settings (key, value) VALUES ('teams_feature_enabled', 'true')
ON CONFLICT (key) DO NOTHING;
