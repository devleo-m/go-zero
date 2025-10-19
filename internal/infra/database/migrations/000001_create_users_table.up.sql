-- Migration: Criar tabela de usuários
-- Versão: 000001
-- Data: 2025-10-19
-- Responsável: devleo-m

-- enum para os papéis do usuário
CREATE TYPE user_role AS ENUM ('admin', 'client');

-- tabela de usuários
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'client',
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Comentários para documentação
COMMENT ON TABLE users IS 'Usuários do sistema GO ZERO';
COMMENT ON COLUMN users.id IS 'Identificador único do usuário (UUID)';
COMMENT ON COLUMN users.email IS 'Email único do usuário (usado para login)';
COMMENT ON COLUMN users.password_hash IS 'Hash da senha do usuário (bcrypt)';
COMMENT ON COLUMN users.full_name IS 'Nome completo do usuário';
COMMENT ON COLUMN users.role IS 'Papel do usuário: admin, client';
COMMENT ON COLUMN users.is_active IS 'Se o usuário está ativo no sistema';
COMMENT ON COLUMN users.is_verified IS 'Se o email foi verificado';
COMMENT ON COLUMN users.created_at IS 'Data de criação do usuário';
COMMENT ON COLUMN users.updated_at IS 'Data da última atualização';
COMMENT ON COLUMN users.deleted_at IS 'Data de exclusão (soft delete)';
