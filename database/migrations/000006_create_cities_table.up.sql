-- Migration: Create Cities Table
-- Description: Create cities lookup table for user location
-- Author: GO ZERO Project

-- Create cities table
CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    state_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_cities_state FOREIGN KEY (state_id) REFERENCES states(id) ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX idx_cities_name ON cities(name);
CREATE INDEX idx_cities_state_id ON cities(state_id);
CREATE INDEX idx_cities_state_name ON cities(state_id, name); -- Composite for performance

-- Insert major cities for São Paulo (state_id = 1)
INSERT INTO cities (name, state_id) VALUES 
('São Paulo', 1),
('Campinas', 1),
('Santos', 1),
('São Bernardo do Campo', 1),
('Santo André', 1),
('Osasco', 1),
('Ribeirão Preto', 1),
('Sorocaba', 1),
('Mauá', 1),
('São José dos Campos', 1);

-- Insert major cities for Rio de Janeiro (state_id = 2)
INSERT INTO cities (name, state_id) VALUES 
('Rio de Janeiro', 2),
('Niterói', 2),
('Nova Iguaçu', 2),
('Duque de Caxias', 2),
('São Gonçalo', 2),
('Campos dos Goytacazes', 2),
('Petrópolis', 2),
('Volta Redonda', 2),
('Magé', 2),
('Cabo Frio', 2);

-- Insert major cities for Minas Gerais (state_id = 3)
INSERT INTO cities (name, state_id) VALUES 
('Belo Horizonte', 3),
('Uberlândia', 3),
('Contagem', 3),
('Juiz de Fora', 3),
('Betim', 3),
('Montes Claros', 3),
('Ribeirão das Neves', 3),
('Uberaba', 3),
('Governador Valadares', 3),
('Ipatinga', 3);

-- Insert major cities for California (state_id = 16)
INSERT INTO cities (name, state_id) VALUES 
('Los Angeles', 16),
('San Francisco', 16),
('San Diego', 16),
('San Jose', 16),
('Fresno', 16),
('Sacramento', 16),
('Long Beach', 16),
('Oakland', 16),
('Bakersfield', 16),
('Anaheim', 16);

-- Insert major cities for Texas (state_id = 17)
INSERT INTO cities (name, state_id) VALUES 
('Houston', 17),
('Dallas', 17),
('Austin', 17),
('San Antonio', 17),
('Fort Worth', 17),
('El Paso', 17),
('Arlington', 17),
('Corpus Christi', 17),
('Plano', 17),
('Lubbock', 17);
