-- Seed: Migrate Users Data
-- Description: Migrate existing users data to use foreign keys
-- Author: GO ZERO Project

-- ==============================================
-- 1. MIGRATE ROLE DATA
-- ==============================================

-- Migrate role data from old columns to new foreign keys
UPDATE users SET role_id = (
    SELECT id FROM roles WHERE name = 'user'
) WHERE role = 'user' AND role_id IS NULL;

UPDATE users SET role_id = (
    SELECT id FROM roles WHERE name = 'admin'  
) WHERE role = 'admin' AND role_id IS NULL;

UPDATE users SET role_id = (
    SELECT id FROM roles WHERE name = 'moderator'
) WHERE role = 'moderator' AND role_id IS NULL;

UPDATE users SET role_id = (
    SELECT id FROM roles WHERE name = 'super_admin'
) WHERE role = 'super_admin' AND role_id IS NULL;

-- Set default role for any unmapped roles
UPDATE users SET role_id = (
    SELECT id FROM roles WHERE name = 'user'
) WHERE role_id IS NULL;

-- ==============================================
-- 2. MIGRATE STATUS DATA
-- ==============================================

-- Migrate status data from old columns to new foreign keys
UPDATE users SET status_id = (
    SELECT id FROM statuses WHERE name = 'active'
) WHERE status = 'active' AND status_id IS NULL;

UPDATE users SET status_id = (
    SELECT id FROM statuses WHERE name = 'inactive'
) WHERE status = 'inactive' AND status_id IS NULL;

UPDATE users SET status_id = (
    SELECT id FROM statuses WHERE name = 'pending'
) WHERE status = 'pending' AND status_id IS NULL;

UPDATE users SET status_id = (
    SELECT id FROM statuses WHERE name = 'suspended'
) WHERE status = 'suspended' AND status_id IS NULL;

-- Set default status for any unmapped statuses
UPDATE users SET status_id = (
    SELECT id FROM statuses WHERE name = 'active'
) WHERE status_id IS NULL;

-- ==============================================
-- 3. MAKE COLUMNS NOT NULL
-- ==============================================

-- Make role_id and status_id NOT NULL after migration
ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;
ALTER TABLE users ALTER COLUMN status_id SET NOT NULL;

-- ==============================================
-- 4. DROP OLD COLUMNS
-- ==============================================

-- Drop old columns after successful migration
ALTER TABLE users DROP COLUMN IF EXISTS role;
ALTER TABLE users DROP COLUMN IF EXISTS status;
