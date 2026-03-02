-- =====================================================
-- THE SENTINEL SYSTEM - INITIAL SCHEMA MIGRATION (DOWN)
-- =====================================================
-- Purpose: Rollback initial schema migration
-- Compatible with: Komgrip Starter Kit
-- Author: Senior Database Architect
-- Date: 2026-01-26
-- =====================================================
-- ⚠️  WARNING: This will permanently delete all data in these tables!
-- ⚠️  Only run this in development or if you need to completely reset!
-- =====================================================

-- ========================================
-- 1. DROP TRIGGERS FIRST
-- ========================================
-- Remove triggers before dropping tables to avoid orphaned triggers
DROP TRIGGER IF EXISTS update_appeals_updated_at ON appeals;
DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;

-- Drop the trigger function (reusable, so only one to drop)
DROP FUNCTION IF EXISTS update_updated_at_column();

-- ========================================
-- 2. DROP TABLES IN REVERSE DEPENDENCY ORDER
-- ========================================
-- Drop in reverse order to avoid foreign key constraint violations
-- Order: child tables first, parent tables last

-- Drop audit_logs (no dependencies, can be dropped anytime)
DROP TABLE IF EXISTS audit_logs CASCADE;

-- Drop appeals (depends on: tasks, users)
DROP TABLE IF EXISTS appeals CASCADE;

-- Drop submissions (depends on: tasks, users)
DROP TABLE IF EXISTS submissions CASCADE;

-- Drop tasks (depends on: users)
DROP TABLE IF EXISTS tasks CASCADE;

-- ========================================
-- 3. REVERT USERS TABLE ALTERATIONS
-- ========================================
-- Remove Sentinel-specific columns from users table
-- Restore to original Komgrip Starter Kit structure

-- Drop indexes first (must drop before dropping columns)
DROP INDEX IF EXISTS idx_users_health_score;
DROP INDEX IF EXISTS idx_users_role;

-- Drop constraints (must drop before dropping columns)
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_health_score;
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;

-- Drop columns (restore to original schema)
ALTER TABLE users DROP COLUMN IF EXISTS tech_stack;
ALTER TABLE users DROP COLUMN IF EXISTS health_score;
ALTER TABLE users DROP COLUMN IF EXISTS role;

-- ========================================
-- 4. OPTIONAL: DROP UUID EXTENSION
-- ========================================
-- ⚠️  CAUTION: Do NOT uncomment unless you're certain no other tables use UUIDs!
-- 
-- The uuid-ossp extension may be used by other parts of the system.
-- Only drop if you're rolling back a complete fresh installation.
-- 
-- Uncomment ONLY if you want to remove UUID support entirely:
-- DROP EXTENSION IF EXISTS "uuid-ossp";

-- =====================================================
-- ROLLBACK COMPLETE ✅
-- =====================================================
-- Summary:
--   - All Sentinel-specific tables dropped (tasks, submissions, appeals, audit_logs)
--   - Users table restored to original Komgrip Starter Kit structure
--   - All triggers and functions cleaned up
--   - Database state reverted to pre-Sentinel installation
-- 
-- The system is now ready for a fresh migration or different schema.
-- =====================================================
