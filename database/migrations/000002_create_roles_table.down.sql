-- Migration Rollback: Drop Roles Table
-- Description: Removes the roles table
-- Author: GO ZERO Project

-- Drop table (this will also drop indexes)
DROP TABLE IF EXISTS roles CASCADE;
