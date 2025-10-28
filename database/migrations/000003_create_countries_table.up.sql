-- Migration: Create Countries Table
-- Description: Create countries lookup table for user location
-- Author: devleo-m

-- Create countries table
CREATE TABLE countries (
    -- Primary key
    id SERIAL PRIMARY KEY,
    
    -- Required fields
    name VARCHAR(100) NOT NULL,
    code VARCHAR(3) NOT NULL UNIQUE, -- ISO 3166-1 alpha-3
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX idx_countries_name ON countries(name);
CREATE INDEX idx_countries_code ON countries(code);
CREATE INDEX idx_countries_deleted_at ON countries(deleted_at);

