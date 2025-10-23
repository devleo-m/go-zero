-- Migration: Create Users Table
-- Description: Creates the users table with all necessary fields for the user domain entity
-- Author: GO ZERO Project
-- Date: 2024

-- Create enum types for role and status
CREATE TYPE user_role AS ENUM ('admin', 'manager', 'user', 'guest');
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'pending', 'suspended');

-- Create users table
CREATE TABLE users (
    -- Base Entity fields (from shared.BaseEntity)
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    
    -- User specific fields
    name VARCHAR(100) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL, -- bcrypt hash
    
    -- Optional fields
    phone VARCHAR(20) NULL,
    
    -- Control fields
    role user_role NOT NULL DEFAULT 'user',
    status user_status NOT NULL DEFAULT 'pending',
    
    -- Login tracking
    last_login_at TIMESTAMP WITH TIME ZONE NULL,
    login_count INTEGER NOT NULL DEFAULT 0,
    
    -- Constraints
    CONSTRAINT users_name_length CHECK (LENGTH(name) >= 2),
    CONSTRAINT users_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT users_login_count_positive CHECK (login_count >= 0)
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at);

-- Create function to automatically update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE users IS 'Users table storing user information and authentication data';
COMMENT ON COLUMN users.id IS 'Unique identifier for the user';
COMMENT ON COLUMN users.name IS 'User full name (2-100 characters)';
COMMENT ON COLUMN users.email IS 'User email address (unique, validated format)';
COMMENT ON COLUMN users.password IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.phone IS 'Optional phone number';
COMMENT ON COLUMN users.role IS 'User role determining permissions (admin, manager, user, guest)';
COMMENT ON COLUMN users.status IS 'User status (active, inactive, pending, suspended)';
COMMENT ON COLUMN users.last_login_at IS 'Timestamp of last successful login';
COMMENT ON COLUMN users.login_count IS 'Number of successful logins';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp (NULL = active, NOT NULL = deleted)';
