-- Migration: Update Users Table with Foreign Keys
-- Description: Add foreign key columns to users table
-- Author: GO ZERO Project

-- ==============================================
-- 1. ADD NEW FOREIGN KEY COLUMNS
-- ==============================================

-- Add new foreign key columns
ALTER TABLE users ADD COLUMN role_id INTEGER;
ALTER TABLE users ADD COLUMN status_id INTEGER;  
ALTER TABLE users ADD COLUMN city_id INTEGER;

-- ==============================================
-- 2. ADD CONSTRAINTS
-- ==============================================

-- Add foreign key constraints
ALTER TABLE users ADD CONSTRAINT fk_users_role 
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT;

ALTER TABLE users ADD CONSTRAINT fk_users_status 
    FOREIGN KEY (status_id) REFERENCES statuses(id) ON DELETE RESTRICT;

ALTER TABLE users ADD CONSTRAINT fk_users_city 
    FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE SET NULL;

-- ==============================================
-- 3. CREATE PERFORMANCE INDEXES
-- ==============================================

-- User indexes for performance
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_status_id ON users(status_id);
CREATE INDEX idx_users_city_id ON users(city_id);

-- Composite indexes for common queries
CREATE INDEX idx_users_status_deleted ON users(status_id, deleted_at);
CREATE INDEX idx_users_role_status ON users(role_id, status_id);
CREATE INDEX idx_users_city_status ON users(city_id, status_id) WHERE city_id IS NOT NULL;
