-- Migration Rollback: Drop States Table
-- Description: Removes the states table
-- Author: GO ZERO Project

-- Drop table (this will also drop indexes and foreign key constraints)
DROP TABLE IF EXISTS states CASCADE;

