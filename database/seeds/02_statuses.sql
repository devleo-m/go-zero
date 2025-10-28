-- Seed: Statuses
-- Description: Insert initial statuses data
-- Author: devleo-m

-- Insert statuses (only if they don't exist)
INSERT INTO statuses (name, description) VALUES 
('active', 'User account is active and functional'),
('inactive', 'User account is temporarily disabled'),
('pending', 'User account is pending approval'),
('suspended', 'User account is suspended due to violations')
ON CONFLICT (name) DO NOTHING;

