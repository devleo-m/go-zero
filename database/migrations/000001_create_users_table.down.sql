-- Migration Rollback: Drop Users Table
-- Description: Removes the users table and related objects
-- Author: GO ZERO Project
-- Date: 2024

-- Drop trigger first
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop table (this will also drop indexes)
DROP TABLE IF EXISTS users;

-- Drop enum types
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;
