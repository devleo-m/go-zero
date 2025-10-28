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

-- Insert major cities for California (state_id = 28)
INSERT INTO cities (name, state_id) VALUES 
('Los Angeles', 28),
('San Francisco', 28),
('San Diego', 28),
('San Jose', 28),
('Fresno', 28),
('Sacramento', 28),
('Long Beach', 28),
('Oakland', 28),
('Bakersfield', 28),
('Anaheim', 28),
('Santa Ana', 28),
('Riverside', 28),
('Stockton', 28),
('Irvine', 28),
('Chula Vista', 28)
ON CONFLICT DO NOTHING;

-- Insert major cities for Texas (state_id = 29)
INSERT INTO cities (name, state_id) VALUES 
('Houston', 29),
('Dallas', 29),
('Austin', 29),
('San Antonio', 29),
('Fort Worth', 29),
('El Paso', 29),
('Arlington', 29),
('Corpus Christi', 29),
('Plano', 29),
('Lubbock', 29),
('Laredo', 29),
('Garland', 29),
('Irving', 29),
('Amarillo', 29)
ON CONFLICT DO NOTHING;

-- Insert major cities for Ontario (state_id = 48)
INSERT INTO cities (name, state_id) VALUES 
('Toronto', 48),
('Ottawa', 48),
('Mississauga', 48),
('Brampton', 48),
('Hamilton', 48),
('London', 48),
('Markham', 48),
('Vaughan', 48),
('Kitchener', 48),
('Windsor', 48),
('Richmond Hill', 48),
('Oakville', 48),
('Burlington', 48),
('Oshawa', 48),
('Barrie', 48)
ON CONFLICT DO NOTHING;

