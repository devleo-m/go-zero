-- Migration Rollback: Drop Roles Table
-- Description: Removes the roles table
-- Author: devleo-m

-- Drop table (this will also drop indexes)
DROP TABLE IF EXISTS roles CASCADE;

