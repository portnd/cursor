-- Revert task story points back to integer values.
ALTER TABLE tasks
    ALTER COLUMN story_points TYPE INTEGER
    USING ROUND(story_points)::INTEGER;
