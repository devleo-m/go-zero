-- Migration Rollback: Drop Countries Table
-- Description: Removes the countries table
-- Author: devleo-m

-- Drop table (this will also drop indexes)
DROP TABLE IF EXISTS countries CASCADE;

