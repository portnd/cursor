-- =====================================================
-- THE SENTINEL SYSTEM - INITIAL SCHEMA MIGRATION (UP)
-- =====================================================
-- Purpose: Initialize core database schema for The Sentinel
-- Compatible with: Komgrip Starter Kit (GORM with INT user IDs)
-- Author: Senior Database Architect
-- Date: 2026-01-26
-- =====================================================

-- ========================================
-- 1. ENABLE UUID EXTENSION
-- ========================================
-- Enable UUID generation functions (uuid_generate_v4)
-- Required for all UUID primary keys in new tables
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ========================================
-- 2. ALTER USERS TABLE
-- ========================================
-- Assumption: users table exists with the following structure:
--   - id SERIAL PRIMARY KEY (auto-increment integer)
--   - email VARCHAR UNIQUE NOT NULL
--   - password VARCHAR NOT NULL
--   - created_at, updated_at TIMESTAMP
-- 
-- Add Sentinel-specific columns

-- Add role column with CHECK constraint
-- Roles: CEO (can assign tasks, review appeals), PM (can assign tasks), DEV (default)
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'DEV';

ALTER TABLE users 
ADD CONSTRAINT check_user_role CHECK (role IN ('CEO', 'PM', 'DEV'));

-- Add health score column (0.00 - 100.00)
-- Tracks developer performance, affected by task completions and rejections
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS health_score DECIMAL(5,2) DEFAULT 100.00;

ALTER TABLE users 
ADD CONSTRAINT check_health_score CHECK (health_score >= 0 AND health_score <= 100);

-- Add tech stack array (PostgreSQL array of text)
-- Example: ['Go', 'PostgreSQL', 'React', 'TypeScript']
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS tech_stack TEXT[];

-- Add index for role filtering (for CEO/PM dashboards)
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Add index for health score sorting
CREATE INDEX IF NOT EXISTS idx_users_health_score ON users(health_score DESC);

-- ========================================
-- 3. CREATE TASKS TABLE
-- ========================================
-- Core task management table
-- Stores all tasks assigned to developers by CEO/PM
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Resource URLs stored as JSONB for flexible structure
    -- Example: {"figma": "https://figma.com/...", "docs": ["https://...", "https://..."], "api_spec": "https://..."}
    resource_urls JSONB DEFAULT '{}',
    
    -- AI-estimated effort in minutes (calculated from task complexity)
    ai_estimated_minutes INT,
    
    -- Task lifecycle timestamps
    due_at TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    
    -- Task status: PENDING, IN_PROGRESS, SUBMITTED, COMPLETED, REJECTED, CANCELLED
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    
    -- Foreign keys to users table (INT type to match existing schema)
    assigned_to INT REFERENCES users(id) ON DELETE SET NULL,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    
    -- Audit timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Business constraints
    CONSTRAINT check_task_status CHECK (status IN ('PENDING', 'IN_PROGRESS', 'SUBMITTED', 'COMPLETED', 'REJECTED', 'CANCELLED')),
    CONSTRAINT check_task_dates CHECK (
        (started_at IS NULL OR started_at >= created_at) AND
        (completed_at IS NULL OR (started_at IS NOT NULL AND completed_at >= started_at))
    ),
    CONSTRAINT check_ai_estimated_minutes CHECK (ai_estimated_minutes IS NULL OR ai_estimated_minutes > 0)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX IF NOT EXISTS idx_tasks_created_by ON tasks(created_by);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_at ON tasks(due_at);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC);

-- Composite index for developer dashboard (assigned tasks filtered by status)
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_status ON tasks(assigned_to, status) WHERE assigned_to IS NOT NULL;

-- Composite index for overdue tasks
CREATE INDEX IF NOT EXISTS idx_tasks_overdue ON tasks(status, due_at) WHERE status NOT IN ('COMPLETED', 'CANCELLED');

-- ========================================
-- 4. CREATE SUBMISSIONS TABLE
-- ========================================
-- Stores code submissions for tasks
-- AI evaluates each submission and provides automated verdict
CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Foreign keys
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    dev_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Git commit information (SHA-1 or SHA-256 hash)
    commit_hash VARCHAR(64) NOT NULL,
    
    -- AI evaluation results
    -- Verdict: PASS (auto-accept), FAIL (reject, can appeal), PENDING (under review)
    ai_verdict VARCHAR(20),
    
    -- AI score (0-100) based on code quality, tests, standards compliance
    ai_score INT,
    
    -- AI feedback stored as JSONB for structured data
    -- Example: {
    --   "code_quality": 85,
    --   "test_coverage": 90,
    --   "security_score": 95,
    --   "performance_score": 80,
    --   "issues": ["Missing error handling in line 45", "Unused import"],
    --   "strengths": ["Good test coverage", "Clean code structure"]
    -- }
    ai_feedback JSONB DEFAULT '{}',
    
    -- Audit timestamp
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Business constraints
    CONSTRAINT check_ai_verdict CHECK (ai_verdict IN ('PASS', 'FAIL', 'PENDING')),
    CONSTRAINT check_ai_score CHECK (ai_score IS NULL OR (ai_score >= 0 AND ai_score <= 100)),
    CONSTRAINT check_commit_hash CHECK (LENGTH(commit_hash) >= 7 AND LENGTH(commit_hash) <= 64)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_submissions_task_id ON submissions(task_id);
CREATE INDEX IF NOT EXISTS idx_submissions_dev_id ON submissions(dev_id);
CREATE INDEX IF NOT EXISTS idx_submissions_created_at ON submissions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_submissions_ai_verdict ON submissions(ai_verdict);

-- Composite index for task submission history (ordered by time)
CREATE INDEX IF NOT EXISTS idx_submissions_task_created ON submissions(task_id, created_at DESC);

-- Composite index for developer submission history
CREATE INDEX IF NOT EXISTS idx_submissions_dev_created ON submissions(dev_id, created_at DESC);

-- Unique constraint: one submission per commit hash per task
CREATE UNIQUE INDEX IF NOT EXISTS idx_submissions_task_commit ON submissions(task_id, commit_hash);

-- ========================================
-- 5. CREATE APPEALS TABLE
-- ========================================
-- Developers can appeal AI FAIL verdicts
-- CEO/PM reviews appeals and makes final human decision
CREATE TABLE IF NOT EXISTS appeals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Foreign keys
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    dev_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Appeal details
    reason TEXT NOT NULL,
    
    -- Appeal status: PENDING (awaiting review), APPROVED (human override), REJECTED (AI was correct)
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    
    -- Admin review (CEO/PM)
    reviewed_by INT REFERENCES users(id) ON DELETE SET NULL,
    admin_comment TEXT,
    
    -- Audit timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Business constraints
    CONSTRAINT check_appeal_status CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
    CONSTRAINT check_reason_length CHECK (LENGTH(reason) >= 10)
);

-- Indexes for common queries (appeals may already exist with submission_id/developer_id from later migrations)
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='appeals' AND column_name='task_id') THEN
    CREATE INDEX IF NOT EXISTS idx_appeals_task_id ON appeals(task_id);
  END IF;
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='appeals' AND column_name='dev_id') THEN
    CREATE INDEX IF NOT EXISTS idx_appeals_dev_id ON appeals(dev_id);
  END IF;
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='appeals' AND column_name='reviewed_by') THEN
    CREATE INDEX IF NOT EXISTS idx_appeals_reviewed_by ON appeals(reviewed_by);
  END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_appeals_status ON appeals(status);
CREATE INDEX IF NOT EXISTS idx_appeals_created_at ON appeals(created_at DESC);

-- Composite index for pending appeals dashboard (CEO/PM view)
CREATE INDEX IF NOT EXISTS idx_appeals_status_created ON appeals(status, created_at DESC) WHERE status = 'PENDING';

-- ========================================
-- 6. CREATE AUDIT_LOGS TABLE
-- ========================================
-- Immutable audit trail for all critical system events
-- Stores: task assignments, submissions, appeals, role changes, health score updates, etc.
-- WHY: Compliance, debugging, fraud detection, system analytics
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Event classification (for filtering and analytics)
    -- Example types:
    --   - TASK_CREATED, TASK_ASSIGNED, TASK_STARTED, TASK_COMPLETED
    --   - SUBMISSION_CREATED, SUBMISSION_AI_EVALUATED
    --   - APPEAL_FILED, APPEAL_REVIEWED
    --   - USER_ROLE_CHANGED, USER_HEALTH_UPDATED
    --   - WALLET_TRANSACTION (future integration)
    event_type VARCHAR(50) NOT NULL,
    
    -- Flexible metadata storage as JSONB
    -- Example for TASK_ASSIGNED: {
    --   "task_id": "uuid-here",
    --   "assigned_to": 123,
    --   "assigned_by": 1,
    --   "task_title": "Implement user authentication",
    --   "estimated_minutes": 240
    -- }
    -- Example for APPEAL_REVIEWED: {
    --   "appeal_id": "uuid-here",
    --   "task_id": "uuid-here",
    --   "dev_id": 123,
    --   "reviewed_by": 1,
    --   "old_status": "PENDING",
    --   "new_status": "APPROVED",
    --   "admin_comment": "Valid concern, AI was too strict"
    -- }
    metadata JSONB DEFAULT '{}',
    
    -- Immutable timestamp (never updated)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for audit queries and analytics
CREATE INDEX IF NOT EXISTS idx_audit_logs_event_type ON audit_logs(event_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- GIN index for JSONB metadata queries (e.g., find all events for a specific task_id)
CREATE INDEX IF NOT EXISTS idx_audit_logs_metadata ON audit_logs USING GIN(metadata);

-- Composite index for event type analytics over time
CREATE INDEX IF NOT EXISTS idx_audit_logs_type_time ON audit_logs(event_type, created_at DESC);

-- ========================================
-- TRIGGERS FOR AUTO-UPDATING updated_at
-- ========================================
-- Automatically update updated_at timestamp on row modifications

-- Create reusable trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to tasks table
DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;
CREATE TRIGGER update_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to appeals table
DROP TRIGGER IF EXISTS update_appeals_updated_at ON appeals;
CREATE TRIGGER update_appeals_updated_at
    BEFORE UPDATE ON appeals
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- MIGRATION COMPLETE ✅
-- =====================================================
-- Summary:
--   - UUID extension enabled
--   - Users table enhanced with role, health_score, tech_stack
--   - 4 new tables created: tasks, submissions, appeals, audit_logs
--   - 25+ indexes for optimal query performance
--   - Auto-update triggers for updated_at columns
--   - All constraints and foreign keys properly configured
-- 
-- Ready for The Sentinel system deployment!
-- =====================================================
