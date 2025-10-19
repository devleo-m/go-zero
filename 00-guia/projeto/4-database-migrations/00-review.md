# 🗄️ ETAPA 4: DATABASE MIGRATIONS

## 🎯 O QUE VAMOS APRENDER NESTA ETAPA

Nesta etapa, implementaremos um **sistema profissional de migrações de banco de dados** que é essencial para qualquer aplicação em produção. É como ter um "controle de versão" para seu banco de dados!

## 🤔 O QUE SÃO MIGRATIONS?

### **DEFINIÇÃO SIMPLES:**
Migrations são **scripts SQL** que descrevem mudanças no banco de dados de forma **versionada** e **reversível**.

### **ANALOGIA:**
Pense como **Git para banco de dados**:
- **Git:** Controla versões do código
- **Migrations:** Controlam versões do banco

```
Código:     v1.0 → v1.1 → v1.2 → v1.3
Banco:      v1.0 → v1.1 → v1.2 → v1.3
```

## 🏗️ ONDE ISSO SE ENCAIXA NA ARQUITETURA?

```
┌─────────────────────────────────────┐
│        INFRASTRUCTURE LAYER         │
│  ┌─────────────┐  ┌─────────────┐   │
│  │  MIGRATIONS │  │  DATABASE   │   │
│  │             │  │             │   │
│  │ • Up/Down   │  │ • GORM      │   │
│  │ • Version   │  │ • Models    │   │
│  │ • History   │  │ • Repos     │   │
│  └─────────────┘  └─────────────┘   │
└─────────────────────────────────────┘
```

## 🎓 CONCEITOS FUNDAMENTAIS

### **1. VERSIONAMENTO**
Cada migration tem um **número de versão** único:
```
001_create_users_table.up.sql
001_create_users_table.down.sql
002_add_email_to_users.up.sql
002_add_email_to_users.down.sql
```

### **2. UP vs DOWN**
- **UP:** Aplica a mudança (cria tabela, adiciona coluna)
- **DOWN:** Reverte a mudança (remove tabela, remove coluna)

### **3. HISTÓRICO**
O banco mantém um registro de quais migrations foram aplicadas:
```sql
CREATE TABLE schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT NOW()
);
```

## 🔧 TECNOLOGIAS QUE VAMOS USAR

### **1. golang-migrate**
**O QUE É:** Ferramenta oficial para migrations em Go

**ANALOGIA:** Como um "gerente de construção" que:
- Lê os projetos (migrations)
- Aplica as mudanças (UP)
- Desfaz se necessário (DOWN)
- Mantém histórico do que foi feito

**CARACTERÍSTICAS:**
- ✅ Suporte a PostgreSQL, MySQL, SQLite
- ✅ CLI para gerenciar migrations
- ✅ Biblioteca Go para usar no código
- ✅ Transações automáticas
- ✅ Rollback seguro

### **2. Estrutura de Arquivos**
```
migrations/
├── 001_create_users_table.up.sql
├── 001_create_users_table.down.sql
├── 002_add_email_to_users.up.sql
├── 002_add_email_to_users.down.sql
├── 003_create_products_table.up.sql
└── 003_create_products_table.down.sql
```

## 🎯 PROBLEMAS QUE RESOLVEMOS

### ❌ **ANTES (Sem migrations):**

```sql
-- Desenvolvedor A cria tabela
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);

-- Desenvolvedor B adiciona coluna
ALTER TABLE users ADD COLUMN email VARCHAR(255);

-- Desenvolvedor C adiciona outra coluna
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
```

**PROBLEMAS:**
- ❌ Cada dev executa SQL diferente
- ❌ Não há histórico do que foi feito
- ❌ Impossível reverter mudanças
- ❌ Deploy manual e propenso a erros
- ❌ Ambientes diferentes ficam desatualizados

### ✅ **DEPOIS (Com migrations):**

```sql
-- 001_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 001_create_users_table.down.sql
DROP TABLE users;

-- 002_add_email_to_users.up.sql
ALTER TABLE users ADD COLUMN email VARCHAR(255) UNIQUE;

-- 002_add_email_to_users.down.sql
ALTER TABLE users DROP COLUMN email;
```

**BENEFÍCIOS:**
- ✅ Histórico completo de mudanças
- ✅ Reversão segura (rollback)
- ✅ Deploy automatizado
- ✅ Ambientes sempre sincronizados
- ✅ Colaboração entre desenvolvedores

## 🚀 FLUXO DE TRABALHO COM MIGRATIONS

### **1. Desenvolvimento Local**
```bash
# Criar nova migration
migrate create -ext sql -dir migrations -seq add_phone_to_users

# Aplicar migrations
migrate -path migrations -database "postgres://..." up

# Reverter última migration
migrate -path migrations -database "postgres://..." down 1

# Ver status
migrate -path migrations -database "postgres://..." version
```

### **2. Deploy em Produção**
```bash
# Aplicar todas as migrations pendentes
migrate -path migrations -database "$DATABASE_URL" up

# Em caso de problema, reverter
migrate -path migrations -database "$DATABASE_URL" down 1
```

### **3. Integração com Go**
```go
// Aplicar migrations automaticamente na inicialização
func runMigrations(databaseURL string) error {
    m, err := migrate.New("file://migrations", databaseURL)
    if err != nil {
        return err
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    
    return nil
}
```

## 📋 TIPOS DE MIGRATIONS

### **1. Estruturais (DDL)**
```sql
-- Criar tabela
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Adicionar coluna
ALTER TABLE users ADD COLUMN email VARCHAR(255);

-- Criar índice
CREATE INDEX idx_users_email ON users(email);

-- Criar constraint
ALTER TABLE users ADD CONSTRAINT uk_users_email UNIQUE(email);
```

### **2. Dados (DML)**
```sql
-- Inserir dados iniciais
INSERT INTO users (name, email) VALUES 
('Admin', 'admin@example.com'),
('User', 'user@example.com');

-- Atualizar dados existentes
UPDATE users SET status = 'active' WHERE status IS NULL;

-- Migrar dados
UPDATE users SET full_name = CONCAT(first_name, ' ', last_name);
```

### **3. Limpeza**
```sql
-- Remover dados desnecessários
DELETE FROM logs WHERE created_at < NOW() - INTERVAL '1 year';

-- Remover colunas não utilizadas
ALTER TABLE users DROP COLUMN old_field;
```

## 🎯 BEST PRACTICES

### **1. Nomenclatura**
```
001_create_users_table.up.sql
002_add_email_to_users.up.sql
003_create_products_table.up.sql
004_add_foreign_key_users_products.up.sql
```

### **2. Sempre Testar DOWN**
```sql
-- UP: Adiciona coluna
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- DOWN: Remove coluna (SEMPRE testar!)
ALTER TABLE users DROP COLUMN phone;
```

### **3. Usar Transações**
```sql
-- Migration complexa com transação
BEGIN;

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    total DECIMAL(10,2)
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

COMMIT;
```

### **4. Validações**
```sql
-- Verificar se coluna existe antes de adicionar
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'email'
    ) THEN
        ALTER TABLE users ADD COLUMN email VARCHAR(255);
    END IF;
END $$;
```

## 🔄 INTEGRAÇÃO COM GORM

### **Por que não usar apenas GORM AutoMigrate?**

**❌ GORM AutoMigrate:**
- Só adiciona, nunca remove
- Não mantém histórico
- Não permite rollback
- Não funciona bem em equipe

**✅ Migrations + GORM:**
- Controle total sobre mudanças
- Histórico completo
- Rollback seguro
- Colaboração perfeita

### **Fluxo Híbrido:**
1. **Migrations:** Estrutura do banco
2. **GORM:** Operações CRUD
3. **Migrations:** Mudanças estruturais

## 🎓 CONCEITOS AVANÇADOS

### **1. Migrations Condicionais**
```sql
-- Só executa se condição for verdadeira
DO $$
BEGIN
    IF (SELECT COUNT(*) FROM users) = 0 THEN
        INSERT INTO users (name, email) VALUES ('Admin', 'admin@example.com');
    END IF;
END $$;
```

### **2. Migrations de Dados Complexos**
```sql
-- Migrar dados de uma estrutura para outra
INSERT INTO new_products (id, name, price, category_id)
SELECT 
    p.id,
    p.name,
    p.price,
    c.id as category_id
FROM old_products p
JOIN categories c ON c.name = p.category_name;
```

### **3. Migrations com Rollback Complexo**
```sql
-- UP: Criar tabela e popular
CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    bio TEXT
);

INSERT INTO user_profiles (user_id, bio)
SELECT id, 'No bio available' FROM users;

-- DOWN: Remover dados e tabela
DELETE FROM user_profiles;
DROP TABLE user_profiles;
```

## 🚨 ARMADILHAS COMUNS

### **1. Não Testar DOWN**
```sql
-- ❌ RUIM: DOWN não funciona
-- UP
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- DOWN (vai quebrar se coluna não existir)
ALTER TABLE users DROP COLUMN phone;

-- ✅ BOM: DOWN seguro
-- UP
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- DOWN
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'phone'
    ) THEN
        ALTER TABLE users DROP COLUMN phone;
    END IF;
END $$;
```

### **2. Migrations Irreversíveis**
```sql
-- ❌ RUIM: Perde dados
DELETE FROM users WHERE status = 'inactive';

-- ✅ BOM: Preserva dados
UPDATE users SET status = 'archived' WHERE status = 'inactive';
```

### **3. Dependências Não Resolvidas**
```sql
-- ❌ RUIM: Ordem errada
-- 001_create_orders.up.sql
CREATE TABLE orders (user_id INTEGER REFERENCES users(id));

-- 002_create_users.up.sql  
CREATE TABLE users (id SERIAL PRIMARY KEY);

-- ✅ BOM: Ordem correta
-- 001_create_users.up.sql
CREATE TABLE users (id SERIAL PRIMARY KEY);

-- 002_create_orders.up.sql
CREATE TABLE orders (user_id INTEGER REFERENCES users(id));
```

## 📊 MONITORAMENTO E OBSERVABILIDADE

### **1. Logs de Migration**
```go
// Log cada migration aplicada
logger.Info("migration_applied",
    zap.String("version", version),
    zap.String("description", description),
    zap.Duration("duration", time.Since(start)),
)
```

### **2. Verificação de Integridade**
```sql
-- Verificar se todas as migrations foram aplicadas
SELECT version FROM schema_migrations ORDER BY version;

-- Verificar integridade do banco
SELECT 
    table_name,
    column_name,
    data_type,
    is_nullable
FROM information_schema.columns 
WHERE table_schema = 'public'
ORDER BY table_name, ordinal_position;
```

## 🎯 EXEMPLOS PRÁTICOS

### **Cenário 1: E-commerce - Criar Tabela de Produtos**
```sql
-- 001_create_products_table.up.sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER DEFAULT 0,
    category_id INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_name ON products(name);

-- 001_create_products_table.down.sql
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_category;
DROP TABLE IF EXISTS products;
```

### **Cenário 2: Adicionar Campo Obrigatório**
```sql
-- 002_add_email_to_users.up.sql
-- Primeiro, popular emails existentes
UPDATE users SET email = CONCAT('user', id, '@example.com') WHERE email IS NULL;

-- Depois, tornar obrigatório
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT uk_users_email UNIQUE(email);

-- 002_add_email_to_users.down.sql
ALTER TABLE users DROP CONSTRAINT IF EXISTS uk_users_email;
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
```

### **Cenário 3: Migração de Dados Complexa**
```sql
-- 003_migrate_user_data.up.sql
-- Criar nova tabela
CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    address TEXT
);

-- Migrar dados
INSERT INTO user_profiles (user_id, first_name, last_name, phone, address)
SELECT 
    id,
    SPLIT_PART(full_name, ' ', 1) as first_name,
    SPLIT_PART(full_name, ' ', 2) as last_name,
    phone,
    address
FROM users;

-- 003_migrate_user_data.down.sql
DROP TABLE IF EXISTS user_profiles;
```

## 🎓 RESUMO DOS CONCEITOS

### **1. Versionamento**
- Cada migration tem número único
- Ordem importa para dependências
- Histórico mantido no banco

### **2. Reversibilidade**
- UP: Aplica mudança
- DOWN: Reverte mudança
- Sempre testar ambos

### **3. Atomicidade**
- Cada migration é uma transação
- Sucesso total ou falha total
- Estado consistente sempre

### **4. Colaboração**
- Mesmo banco para todos
- Deploy automatizado
- Sincronização de ambientes

## 🚀 BENEFÍCIOS ALCANÇADOS

✅ **Controle de versão** do banco de dados
✅ **Deploy automatizado** e confiável
✅ **Rollback seguro** em caso de problemas
✅ **Colaboração perfeita** entre desenvolvedores
✅ **Ambientes sincronizados** (dev, staging, prod)
✅ **Histórico completo** de mudanças
✅ **Integridade de dados** preservada

## 🎯 PRÓXIMOS PASSOS

Agora que entendemos **TUDO** sobre migrations, na próxima etapa vamos:

1. **Instalar** golang-migrate
2. **Criar** estrutura de migrations
3. **Implementar** primeira migration (users)
4. **Integrar** com nossa aplicação
5. **Testar** UP e DOWN
6. **Automatizar** via Makefile

## 🏆 RESUMO DO QUE VAMOS IMPLEMENTAR

✅ **Sistema de migrations** com golang-migrate
✅ **Estrutura de arquivos** organizada
✅ **Integração com Go** na inicialização
✅ **Comandos Makefile** para facilitar uso
✅ **Validação e logs** de migrations
✅ **Rollback seguro** implementado
✅ **Documentação completa** criada

**Resultado:** Um sistema profissional de controle de versão do banco de dados! 🚀
