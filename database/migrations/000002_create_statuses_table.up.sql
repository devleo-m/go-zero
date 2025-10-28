-- Migration: Create Statuses Table
-- Description: Create statuses lookup table for user status
-- Author: GO ZERO Project

-- Create statuses table
CREATE TABLE statuses (
    -- Primary key
    id SERIAL PRIMARY KEY,
    
    -- Required fields
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX idx_statuses_name ON statuses(name);
CREATE INDEX idx_statuses_deleted_at ON statuses(deleted_at);

