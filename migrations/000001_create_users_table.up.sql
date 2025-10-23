-- ==========================================
-- MIGRATION: CREATE USERS TABLE
-- ==========================================
-- Descrição: Cria a tabela principal de usuários
-- Data: 2024-01-01
-- Autor: GO ZERO Team
-- ==========================================

-- Habilitar extensão UUID se não existir
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criar tabela principal de usuários
CREATE TABLE users (
    -- ==========================================
    -- CAMPOS BASE (AUDITORIA)
    -- ==========================================
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    -- ==========================================
    -- CAMPOS ESSENCIAIS DO USUÁRIO
    -- ==========================================
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    role VARCHAR(50) NOT NULL DEFAULT 'user',

    -- ==========================================
    -- CONSTRAINTS
    -- ==========================================
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_status_check CHECK (status IN ('pending', 'active', 'inactive', 'suspended')),
    CONSTRAINT users_role_check CHECK (role IN ('admin', 'manager', 'user', 'guest')),
    CONSTRAINT users_name_length_check CHECK (LENGTH(name) >= 2 AND LENGTH(name) <= 100),
    CONSTRAINT users_email_format_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT users_phone_format_check CHECK (phone IS NULL OR phone ~* '^\+?[1-9]\d{1,14}$')
);

-- ==========================================
-- ÍNDICES PARA PERFORMANCE
-- ==========================================

-- Índice principal por ID (já existe como PRIMARY KEY)
-- CREATE INDEX idx_users_id ON users(id); -- Desnecessário (PK)

-- Índices para campos de busca frequente
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_updated_at ON users(updated_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Índices para campos de filtro
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_email ON users(email);

-- Índice para soft delete
CREATE INDEX idx_users_active ON users(deleted_at) WHERE deleted_at IS NULL;

-- Índice composto para queries comuns
CREATE INDEX idx_users_status_created_at ON users(status, created_at DESC);
CREATE INDEX idx_users_role_status ON users(role, status);

-- ==========================================
-- TRIGGERS PARA AUDITORIA
-- ==========================================

-- Função para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para atualizar updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ==========================================
-- COMENTÁRIOS PARA DOCUMENTAÇÃO
-- ==========================================

COMMENT ON TABLE users IS 'Tabela principal de usuários do sistema';
COMMENT ON COLUMN users.id IS 'Identificador único do usuário (UUID)';
COMMENT ON COLUMN users.created_at IS 'Data e hora de criação do usuário';
COMMENT ON COLUMN users.updated_at IS 'Data e hora da última atualização';
COMMENT ON COLUMN users.deleted_at IS 'Data e hora do soft delete (NULL = ativo)';
COMMENT ON COLUMN users.name IS 'Nome completo do usuário (2-100 caracteres)';
COMMENT ON COLUMN users.email IS 'Email único do usuário (formato válido)';
COMMENT ON COLUMN users.password_hash IS 'Hash da senha (bcrypt)';
COMMENT ON COLUMN users.phone IS 'Telefone do usuário (formato internacional)';
COMMENT ON COLUMN users.status IS 'Status do usuário: pending, active, inactive, suspended';
COMMENT ON COLUMN users.role IS 'Papel do usuário: admin, manager, user, guest';

-- ==========================================
-- DADOS INICIAIS (OPCIONAL)
-- ==========================================

-- Inserir usuário admin padrão (senha: admin123)
-- INSERT INTO users (name, email, password_hash, status, role) VALUES
-- ('Administrador', 'admin@go-zero.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active', 'admin');

-- ==========================================
-- VERIFICAÇÕES FINAIS
-- ==========================================

-- Verificar se a tabela foi criada corretamente
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        RAISE EXCEPTION 'Falha ao criar tabela users';
    END IF;
    
    RAISE NOTICE 'Tabela users criada com sucesso!';
END $$;
