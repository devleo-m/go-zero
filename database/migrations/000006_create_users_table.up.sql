-- Migration: Create Users Table
-- Description: Create users table with all required fields and foreign keys
-- Author: devleo-m

-- Create users table
CREATE TABLE users (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Required fields
    name VARCHAR(100) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    
    -- Optional fields
    phone VARCHAR(20),
    avatar VARCHAR(500),
    bio TEXT,
    data_nascimento DATE,
    cpf VARCHAR(14) UNIQUE,
    
    -- Foreign keys
    role_id INTEGER NOT NULL,
    status_id INTEGER NOT NULL,
    city_id INTEGER,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT,
    CONSTRAINT fk_users_status FOREIGN KEY (status_id) REFERENCES statuses(id) ON DELETE RESTRICT,
    CONSTRAINT fk_users_city FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE SET NULL
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_cpf ON users(cpf);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_status_id ON users(status_id);
CREATE INDEX idx_users_city_id ON users(city_id);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Composite indexes for common queries
CREATE INDEX idx_users_status_deleted ON users(status_id, deleted_at);
CREATE INDEX idx_users_role_status ON users(role_id, status_id);
CREATE INDEX idx_users_city_status ON users(city_id, status_id) WHERE city_id IS NOT NULL;

