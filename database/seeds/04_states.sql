-- Seed: States
-- Description: Insert initial states data
-- Author: GO ZERO Project

-- Insert states for Brazil (country_id = 1)
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
('Piauí', 1, 'PI'),
('Alagoas', 1, 'AL'),
('Tocantins', 1, 'TO'),
('Rio Grande do Norte', 1, 'RN'),
('Acre', 1, 'AC'),
('Amapá', 1, 'AP'),
('Amazonas', 1, 'AM'),
('Mato Grosso', 1, 'MT'),
('Mato Grosso do Sul', 1, 'MS'),
('Rondônia', 1, 'RO'),
('Roraima', 1, 'RR'),
('Sergipe', 1, 'SE'),
('Distrito Federal', 1, 'DF')
ON CONFLICT DO NOTHING;

-- Insert states for USA (country_id = 2)
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
('Michigan', 2, 'MI'),
('New Jersey', 2, 'NJ'),
('Virginia', 2, 'VA'),
('Washington', 2, 'WA'),
('Arizona', 2, 'AZ'),
('Massachusetts', 2, 'MA'),
('Tennessee', 2, 'TN'),
('Indiana', 2, 'IN'),
('Missouri', 2, 'MO'),
('Maryland', 2, 'MD'),
('Wisconsin', 2, 'WI')
ON CONFLICT DO NOTHING;

-- Insert states for Canada (country_id = 3)
INSERT INTO states (name, country_id, code) VALUES 
('Ontario', 3, 'ON'),
('Quebec', 3, 'QC'),
('British Columbia', 3, 'BC'),
('Alberta', 3, 'AB'),
('Manitoba', 3, 'MB'),
('Saskatchewan', 3, 'SK'),
('Nova Scotia', 3, 'NS'),
('New Brunswick', 3, 'NB'),
('Newfoundland and Labrador', 3, 'NL'),
('Prince Edward Island', 3, 'PE')
ON CONFLICT DO NOTHING;
