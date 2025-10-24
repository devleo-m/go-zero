-- Migration: Create Roles Table
-- Description: Create roles lookup table for user roles
-- Author: GO ZERO Project

-- Create roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_roles_name ON roles(name);
CREATE UNIQUE INDEX idx_roles_name_unique ON roles(name);

-- Insert initial roles
INSERT INTO roles (name, description) VALUES 
('user', 'Regular user with basic permissions'),
('admin', 'Administrator with elevated permissions'),
('moderator', 'Moderator with content management permissions'),
('super_admin', 'Super administrator with full system access');
