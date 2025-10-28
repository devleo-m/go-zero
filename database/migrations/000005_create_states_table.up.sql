-- Migration: Create States Table
-- Description: Create states lookup table for user location
-- Author: GO ZERO Project

-- Create states table
CREATE TABLE states (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    country_id INTEGER NOT NULL,
    code VARCHAR(10), -- State/Province code
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_states_country FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX idx_states_name ON states(name);
CREATE INDEX idx_states_country_id ON states(country_id);
CREATE INDEX idx_states_code ON states(code);
CREATE INDEX idx_states_country_name ON states(country_id, name); -- Composite for performance

-- Insert initial states for Brazil
INSERT INTO states (name, country_id, code) VALUES 
('São Paulo', 1, 'SP'),
('Rio de Janeiro', 1, 'RJ'),
('Minas Gerais', 1, 'MG'),
('Bahia', 1, 'BA'),
('Paraná', 1, 'PR'),
('Rio Grande do Sul', 1, 'RS'),
('Pernambuco', 1, 'PE'),
('Ceará', 1, 'CE'),
('Pará', 1, 'PA'),
('Santa Catarina', 1, 'SC'),
('Goiás', 1, 'GO'),
('Maranhão', 1, 'MA'),
('Paraíba', 1, 'PB'),
('Espírito Santo', 1, 'ES'),
('Piauí', 1, 'PI');

-- Insert initial states for USA
INSERT INTO states (name, country_id, code) VALUES 
('California', 2, 'CA'),
('Texas', 2, 'TX'),
('Florida', 2, 'FL'),
('New York', 2, 'NY'),
('Illinois', 2, 'IL'),
('Pennsylvania', 2, 'PA'),
('Ohio', 2, 'OH'),
('Georgia', 2, 'GA'),
('North Carolina', 2, 'NC'),
('Michigan', 2, 'MI');
