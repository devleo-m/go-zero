-- Seed: Roles
-- Description: Insert initial roles data
-- Author: GO ZERO Project

-- Insert roles (only if they don't exist)
INSERT INTO roles (name, description) VALUES 
('user', 'Regular user with basic permissions'),
('admin', 'Administrator with elevated permissions'),
('moderator', 'Moderator with content management permissions'),
('super_admin', 'Super administrator with full system access')
ON CONFLICT (name) DO NOTHING;
