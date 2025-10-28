-- Migration Rollback: Drop Cities Table
-- Description: Removes the cities table
-- Author: devleo-m

-- Drop table (this will also drop indexes and foreign key constraints)
DROP TABLE IF EXISTS cities CASCADE;

