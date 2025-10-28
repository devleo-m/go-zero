-- Migration: Create Statuses Table
-- Description: Create statuses lookup table for user status
-- Author: GO ZERO Project

-- Create statuses table
CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_statuses_name ON statuses(name);
CREATE UNIQUE INDEX idx_statuses_name_unique ON statuses(name);

-- Insert initial statuses
INSERT INTO statuses (name, description) VALUES 
('active', 'User account is active and functional'),
('inactive', 'User account is temporarily disabled'),
('pending', 'User account is pending approval'),
('suspended', 'User account is suspended due to violations');
