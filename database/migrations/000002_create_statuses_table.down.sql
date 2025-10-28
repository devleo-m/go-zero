-- Migration Rollback: Drop Statuses Table
-- Description: Removes the statuses table
-- Author: GO ZERO Project

-- Drop table (this will also drop indexes)
DROP TABLE IF EXISTS statuses CASCADE;

