# 🗄️ COMANDOS DE MIGRATIONS - GO ZERO

## 🎯 **GUIA PRÁTICO DE USO**

Este guia contém **TODOS** os comandos necessários para trabalhar com migrations no projeto GO ZERO.

---

## 📋 **COMANDOS BÁSICOS**

### **1. Aplicar Migrations (UP)**
```bash
# Aplicar TODAS as migrations pendentes
go run cmd/migrate/main.go -direction=up

# Aplicar apenas 1 migration
go run cmd/migrate/main.go -direction=up -steps=1

# Aplicar 3 migrations
go run cmd/migrate/main.go -direction=up -steps=3
```

### **2. Reverter Migrations (DOWN)**
```bash
# Reverter APENAS a última migration
go run cmd/migrate/main.go -direction=down -steps=1

# Reverter 2 migrations
go run cmd/migrate/main.go -direction=down -steps=2

# Reverter TODAS as migrations (CUIDADO!)
go run cmd/migrate/main.go -direction=down
```

### **3. Forçar Versão Específica**
```bash
# Forçar para versão 1 (útil quando migration falha)
go run cmd/migrate/main.go -direction=force -steps=1

# Forçar para versão 0 (limpar tudo)
go run cmd/migrate/main.go -direction=force -steps=0
```

---

## 🛠️ **COMANDOS DE DESENVOLVIMENTO**

### **4. Criar Nova Migration**
```bash
# Navegar para pasta de migrations
cd internal/infra/database/migrations

# Criar migration (substitua 'nome_da_migration' pelo nome desejado)
migrate create -ext sql -dir . -seq nome_da_migration

# Exemplos:
migrate create -ext sql -dir . -seq add_phone_to_users
migrate create -ext sql -dir . -seq create_products_table
migrate create -ext sql -dir . -seq add_email_index
```

### **5. Verificar Status das Migrations**
```bash
# Verificar versão atual (sem aplicar migrations)
go run cmd/migrate/main.go -direction=up -steps=0
```

---

## 🔍 **COMANDOS DE VERIFICAÇÃO**

### **6. Verificar Tabelas no Banco**
```bash
# Listar todas as tabelas
docker exec go-zero-db psql -U postgres -d go_zero_dev -c "\dt"

# Ver estrutura da tabela users
docker exec go-zero-db psql -U postgres -d go_zero_dev -c "\d users"

# Ver índices da tabela users
docker exec go-zero-db psql -U postgres -d go_zero_dev -c "\di users"
```

### **7. Verificar Histórico de Migrations**
```bash
# Ver tabela de controle de migrations
docker exec go-zero-db psql -U postgres -d go_zero_dev -c "SELECT * FROM schema_migrations ORDER BY version;"
```

---

## 🚨 **COMANDOS DE EMERGÊNCIA**

### **8. Reset Completo (CUIDADO!)**
```bash
# 1. Reverter todas as migrations
go run cmd/migrate/main.go -direction=down

# 2. Aplicar novamente
go run cmd/migrate/main.go -direction=up
```

### **9. Limpar Banco Completamente**
```bash
# Conectar ao banco e limpar manualmente
docker exec -it go-zero-db psql -U postgres -d go_zero_dev

# Dentro do psql:
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
\q

# Aplicar migrations novamente
go run cmd/migrate/main.go -direction=up
```

---

## 📝 **FLUXO DE TRABALHO RECOMENDADO**

### **Cenário 1: Desenvolvimento Normal**
```bash
# 1. Criar nova migration
cd internal/infra/database/migrations
migrate create -ext sql -dir . -seq add_user_phone

# 2. Editar arquivos .up.sql e .down.sql
# 3. Aplicar migration
cd ../../..
go run cmd/migrate/main.go -direction=up

# 4. Verificar se funcionou
docker exec go-zero-db psql -U postgres -d go_zero_dev -c "\d users"
```

### **Cenário 2: Testar Rollback**
```bash
# 1. Aplicar migration
go run cmd/migrate/main.go -direction=up

# 2. Testar rollback
go run cmd/migrate/main.go -direction=down -steps=1

# 3. Aplicar novamente
go run cmd/migrate/main.go -direction=up
```

### **Cenário 3: Migration Falhou**
```bash
# 1. Verificar status
go run cmd/migrate/main.go -direction=up -steps=0

# 2. Forçar versão anterior
go run cmd/migrate/main.go -direction=force -steps=1

# 3. Corrigir migration e aplicar novamente
go run cmd/migrate/main.go -direction=up
```

---

## 🎯 **EXEMPLOS PRÁTICOS**

### **Exemplo 1: Adicionar Campo à Tabela Users**
```bash
# 1. Criar migration
cd internal/infra/database/migrations
migrate create -ext sql -dir . -seq add_phone_to_users

# 2. Editar 000002_add_phone_to_users.up.sql:
# ALTER TABLE users ADD COLUMN phone VARCHAR(20);

# 3. Editar 000002_add_phone_to_users.down.sql:
# ALTER TABLE users DROP COLUMN phone;

# 4. Aplicar
cd ../../..
go run cmd/migrate/main.go -direction=up
```

### **Exemplo 2: Criar Nova Tabela**
```bash
# 1. Criar migration
cd internal/infra/database/migrations
migrate create -ext sql -dir . -seq create_products_table

# 2. Editar arquivos SQL
# 3. Aplicar
cd ../../..
go run cmd/migrate/main.go -direction=up
```

---

## ⚠️ **REGRAS IMPORTANTES**

### **✅ FAÇA:**
- Sempre teste UP e DOWN antes de commitar
- Use nomes descritivos para migrations
- Faça backup antes de migrations em produção
- Use transações para migrations complexas

### **❌ NÃO FAÇA:**
- Editar migrations já aplicadas em produção
- Fazer migrations irreversíveis (DELETE sem WHERE)
- Aplicar migrations sem testar
- Usar nomes genéricos (migration1, fix, etc.)

---

## 🆘 **SOLUÇÃO DE PROBLEMAS**

### **Problema: Migration falhou no meio**
```bash
# Solução: Forçar versão anterior
go run cmd/migrate/main.go -direction=force -steps=1
```

### **Problema: Banco está "dirty" (migration incompleta)**
```bash
# Solução: Forçar versão limpa
go run cmd/migrate/main.go -direction=force -steps=0
go run cmd/migrate/main.go -direction=up
```

### **Problema: Migration não aparece**
```bash
# Verificar se arquivos existem
ls internal/infra/database/migrations/

# Verificar se banco está rodando
docker ps | grep postgres
```

---

## 📊 **COMANDOS DE MONITORAMENTO**

### **Verificar Performance**
```bash
# Ver queries lentas
docker exec go-zero-db psql -U postgres -d go_zero_dev -c "
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;"
```

### **Verificar Espaço Usado**
```bash
# Tamanho das tabelas
docker exec go-zero-db psql -U postgres -d go_zero_dev -c "
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
```

---

## 🎉 **RESUMO DOS COMANDOS MAIS USADOS**

| Ação | Comando |
|------|---------|
| **Aplicar migrations** | `go run cmd/migrate/main.go -direction=up` |
| **Reverter última** | `go run cmd/migrate/main.go -direction=down -steps=1` |
| **Criar migration** | `migrate create -ext sql -dir . -seq nome` |
| **Ver tabelas** | `docker exec go-zero-db psql -U postgres -d go_zero_dev -c "\dt"` |
| **Ver estrutura** | `docker exec go-zero-db psql -U postgres -d go_zero_dev -c "\d users"` |
| **Forçar versão** | `go run cmd/migrate/main.go -direction=force -steps=1` |

---

## 🚀 **PRÓXIMOS PASSOS**

Agora que você domina as migrations, pode:

1. **Criar mais tabelas** (products, orders, etc.)
2. **Adicionar campos** às tabelas existentes
3. **Criar índices** para performance
4. **Implementar soft delete** em outras tabelas
5. **Criar views** e procedures

**Lembre-se:** Sempre teste UP e DOWN antes de commitar! 🎯

