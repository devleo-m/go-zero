-- Migration Rollback: Drop Users Table
-- Description: Removes the users table
-- Author: GO ZERO Project

-- Drop table (this will also drop indexes and foreign key constraints)
DROP TABLE IF EXISTS users CASCADE;

