-- ==========================================
-- MIGRATION DOWN: DROP USERS TABLE
-- ==========================================
-- Descrição: Remove a tabela de usuários e todos os objetos relacionados
-- Data: 2024-01-01
-- Autor: GO ZERO Team
-- ==========================================

-- ==========================================
-- REMOVER TRIGGERS
-- ==========================================

-- Remover trigger de updated_at
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- ==========================================
-- REMOVER FUNÇÕES
-- ==========================================

-- Remover função de updated_at (apenas se não for usada por outras tabelas)
-- DROP FUNCTION IF EXISTS update_updated_at_column();

-- ==========================================
-- REMOVER ÍNDICES
-- ==========================================

-- Remover índices criados
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_updated_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_status_created_at;
DROP INDEX IF EXISTS idx_users_role_status;

-- ==========================================
-- REMOVER TABELA PRINCIPAL
-- ==========================================

-- Remover tabela de usuários
DROP TABLE IF EXISTS users CASCADE;

-- ==========================================
-- VERIFICAÇÕES FINAIS
-- ==========================================

-- Verificar se a tabela foi removida corretamente
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        RAISE EXCEPTION 'Falha ao remover tabela users';
    END IF;
    
    RAISE NOTICE 'Tabela users removida com sucesso!';
END $$;
