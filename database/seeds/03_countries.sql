-- Seed: Countries
-- Description: Insert initial countries data
-- Author: devleo-m

-- Insert countries (only if they don't exist)
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
('Paraguay', 'PRY')
ON CONFLICT (code) DO NOTHING;

