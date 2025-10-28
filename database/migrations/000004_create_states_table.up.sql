-- Migration: Create States Table
-- Description: Create states lookup table for user location
-- Author: GO ZERO Project

-- Create states table
CREATE TABLE states (
    -- Primary key
    id SERIAL PRIMARY KEY,
    
    -- Required fields
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10), -- State/Province code
    
    -- Foreign keys
    country_id INTEGER NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_states_country FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX idx_states_name ON states(name);
CREATE INDEX idx_states_code ON states(code);
CREATE INDEX idx_states_country_id ON states(country_id);
CREATE INDEX idx_states_deleted_at ON states(deleted_at);

-- Composite indexes for common queries
CREATE INDEX idx_states_country_name ON states(country_id, name);

