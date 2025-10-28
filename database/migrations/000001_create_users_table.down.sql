-- Migration Rollback: Drop Users Table
-- Description: Removes the users table
-- Author: GO ZERO Project

-- Drop table (this will also drop indexes)
DROP TABLE IF EXISTS users;