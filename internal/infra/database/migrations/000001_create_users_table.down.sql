-- Rollback: Remover tabela de usuários
-- Versão: 000001
-- Data: 2025-10-19
-- Responsável: devleo-m

-- Remover índices primeiro (ordem reversa)
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;

-- Remover tabela
DROP TABLE IF EXISTS users;
