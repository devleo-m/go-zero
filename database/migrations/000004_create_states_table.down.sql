-- Migration Rollback: Drop States Table
-- Description: Removes the states table
-- Author: devleo-m

-- Drop table (this will also drop indexes and foreign key constraints)
DROP TABLE IF EXISTS states CASCADE;

