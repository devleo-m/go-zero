-- Migration Rollback: Drop Countries Table
-- Description: Removes the countries table
-- Author: GO ZERO Project

-- Drop table (this will also drop indexes)
DROP TABLE IF EXISTS countries CASCADE;

