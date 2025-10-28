-- Seed: Cities
-- Description: Insert initial cities data
-- Author: GO ZERO Project

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
('São José dos Campos', 1),
('Mogi das Cruzes', 1),
('Diadema', 1),
('Jundiaí', 1),
('Carapicuíba', 1),
('Piracicaba', 1)
ON CONFLICT DO NOTHING;

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
('Cabo Frio', 2),
('Nova Friburgo', 2),
('Barra Mansa', 2),
('Angra dos Reis', 2),
('Teresópolis', 2),
('Nilópolis', 2)
ON CONFLICT DO NOTHING;

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
('Ipatinga', 3),
('Sete Lagoas', 3),
('Divinópolis', 3),
('Santa Luzia', 3),
('Ibirité', 3),
('Poços de Caldas', 3)
ON CONFLICT DO NOTHING;

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
('Anaheim', 16),
('Santa Ana', 16),
('Riverside', 16),
('Stockton', 16),
('Irvine', 16),
('Chula Vista', 16)
ON CONFLICT DO NOTHING;

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
('Lubbock', 17),
('Laredo', 17),
('Lubbock', 17),
('Garland', 17),
('Irving', 17),
('Amarillo', 17)
ON CONFLICT DO NOTHING;

-- Insert major cities for Ontario (state_id = 26)
INSERT INTO cities (name, state_id) VALUES 
('Toronto', 26),
('Ottawa', 26),
('Mississauga', 26),
('Brampton', 26),
('Hamilton', 26),
('London', 26),
('Markham', 26),
('Vaughan', 26),
('Kitchener', 26),
('Windsor', 26),
('Richmond Hill', 26),
('Oakville', 26),
('Burlington', 26),
('Oshawa', 26),
('Barrie', 26)
ON CONFLICT DO NOTHING;
