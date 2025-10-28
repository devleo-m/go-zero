-- Migration: Create Roles Table
-- Description: Create roles lookup table for user roles
-- Author: GO ZERO Project

-- Create roles table
CREATE TABLE roles (
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
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_deleted_at ON roles(deleted_at);

