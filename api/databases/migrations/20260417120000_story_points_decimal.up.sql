-- Convert task story points to decimal to support half-point estimates.
ALTER TABLE tasks
    ALTER COLUMN story_points TYPE DECIMAL(5,1)
    USING story_points::DECIMAL(5,1);
