-- Migration: Create Cities Table
-- Description: Create cities lookup table for user location
-- Author: devleo-m

-- Create cities table
CREATE TABLE cities (
    -- Primary key
    id SERIAL PRIMARY KEY,
    
    -- Required fields
    name VARCHAR(100) NOT NULL,
    
    -- Foreign keys
    state_id INTEGER NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_cities_state FOREIGN KEY (state_id) REFERENCES states(id) ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX idx_cities_name ON cities(name);
CREATE INDEX idx_cities_state_id ON cities(state_id);
CREATE INDEX idx_cities_deleted_at ON cities(deleted_at);

-- Composite indexes for common queries
CREATE INDEX idx_cities_state_name ON cities(state_id, name);

