-- Migration Rollback: Revert Users Table Foreign Keys
-- Description: Remove foreign key columns from users table
-- Author: GO ZERO Project

-- ==============================================
-- 1. DROP FOREIGN KEY CONSTRAINTS
-- ==============================================

-- Drop foreign key constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_role;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_status;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_city;

-- ==============================================
-- 2. DROP NEW COLUMNS
-- ==============================================

-- Drop new columns
ALTER TABLE users DROP COLUMN IF EXISTS role_id;
ALTER TABLE users DROP COLUMN IF EXISTS status_id;
ALTER TABLE users DROP COLUMN IF EXISTS city_id;

-- ==============================================
-- 3. DROP NEW INDEXES
-- ==============================================

-- Drop new indexes
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_users_status_id;
DROP INDEX IF EXISTS idx_users_city_id;
DROP INDEX IF EXISTS idx_users_status_deleted;
DROP INDEX IF EXISTS idx_users_role_status;
DROP INDEX IF EXISTS idx_users_city_status;
