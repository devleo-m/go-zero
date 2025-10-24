-- Migration: Create Countries Table
-- Description: Create countries lookup table for user location
-- Author: GO ZERO Project

-- Create countries table
CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(3) NOT NULL UNIQUE, -- ISO 3166-1 alpha-3
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_countries_name ON countries(name);
CREATE INDEX idx_countries_code ON countries(code);
CREATE UNIQUE INDEX idx_countries_code_unique ON countries(code);

-- Insert initial countries
INSERT INTO countries (name, code) VALUES 
('Brazil', 'BRA'),
('United States', 'USA'),
('Canada', 'CAN'),
('Argentina', 'ARG'),
('Chile', 'CHL'),
('Mexico', 'MEX'),
('Colombia', 'COL'),
('Peru', 'PER'),
('Uruguay', 'URY'),
('Paraguay', 'PRY');
