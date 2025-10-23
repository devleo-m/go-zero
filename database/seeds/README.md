# Sistema de Seeds - GO ZERO

## O que são Seeds?

Seeds são dados iniciais que populam o banco de dados para desenvolvimento e testes. Diferente das migrations que criam a estrutura do banco, os seeds inserem dados de exemplo.

## Estrutura do Sistema

```
database/
├── seeds/
│   ├── seeder.go           # Gerenciador principal de seeders
│   └── users_seeder.go     # Seeder específico de usuários
cmd/
└── seed/
    └── main.go             # Comando para executar seeds
```

## Como Usar

### Comandos Disponíveis

```bash
# Executar todos os seeds
make seed

# Executar apenas seed de usuários
make seed-users

# Limpar todos os dados seedados
make seed-clear

# Limpar apenas dados de usuários seedados
make seed-clear-users

# Setup completo (migrations + seeds)
make dev-setup
```

### Comandos Diretos

```bash
# Executar todos os seeds
go run cmd/seed/main.go

# Executar apenas seed de usuários
go run cmd/seed/main.go -action=users

# Limpar todos os dados seedados
go run cmd/seed/main.go -action=clear

# Limpar apenas dados de usuários seedados
go run cmd/seed/main.go -action=clear-users

# Ver ajuda
go run cmd/seed/main.go -help
```

## Usuários Criados pelos Seeds

O sistema cria os seguintes usuários para desenvolvimento:

| Nome | Email | Senha | Role | Status |
|------|-------|-------|------|--------|
| Administrador do Sistema | admin@go-zero.com | Admin123!@# | admin | active |
| Gerente de Projetos | manager@go-zero.com | Manager123!@# | manager | active |
| Usuário Teste | user@go-zero.com | User123!@# | user | active |
| Usuário Pendente | pending@go-zero.com | Pending123!@# | user | pending |
| Usuário Suspenso | suspended@go-zero.com | Suspended123!@# | user | suspended |
| Convidado | guest@go-zero.com | Guest123!@# | guest | active |

## Características do Sistema

### ✅ Segurança
- Senhas são hasheadas com bcrypt
- Validação de formato de email
- Validação de força da senha
- Verificação de duplicatas

### ✅ Flexibilidade
- Sistema modular (cada entidade tem seu seeder)
- Comandos específicos por entidade
- Limpeza seletiva de dados

### ✅ Robustez
- Verificação de existência antes de criar
- Tratamento de erros completo
- Logs detalhados de operações

### ✅ Manutenibilidade
- Código organizado em camadas
- Documentação completa
- Fácil adição de novos seeders

## Adicionando Novos Seeders

### 1. Criar o Seeder

```go
// database/seeds/products_seeder.go
package seeds

import (
    "log"
    "gorm.io/gorm"
)

type ProductSeeder struct {
    db *gorm.DB
}

func NewProductSeeder(db *gorm.DB) *ProductSeeder {
    return &ProductSeeder{db: db}
}

func (s *ProductSeeder) SeedProducts() error {
    log.Println("🌱 Starting product seeding...")
    
    // Lógica de seeding aqui
    
    log.Println("✅ Product seeding completed!")
    return nil
}

func (s *ProductSeeder) ClearProducts() error {
    log.Println("🧹 Clearing product seeded data...")
    
    // Lógica de limpeza aqui
    
    log.Println("✅ Product clearing completed!")
    return nil
}
```

### 2. Adicionar ao SeederManager

```go
// database/seeds/seeder.go
func (sm *SeederManager) RunAll() error {
    // ... seeders existentes ...
    
    productSeeder := NewProductSeeder(sm.db)
    if err := productSeeder.SeedProducts(); err != nil {
        return err
    }
    
    return nil
}
```

### 3. Adicionar Comandos ao Makefile

```makefile
seed-products: ## Executa apenas seed de produtos
	@echo "$(BLUE)📦 Executando seed de produtos...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go run cmd/seed/main.go -action=products

seed-clear-products: ## Limpa apenas dados de produtos seedados
	@echo "$(BLUE)🧹 Limpando dados de produtos seedados...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go run cmd/seed/main.go -action=clear-products
```

## Boas Práticas

### ✅ Faça
- Use dados realistas para desenvolvimento
- Sempre verifique se já existe antes de criar
- Use senhas seguras mas fáceis de lembrar para dev
- Documente os dados criados
- Implemente limpeza seletiva

### ❌ Não Faça
- Não use senhas fracas em produção
- Não crie dados duplicados
- Não esqueça de implementar a função de limpeza
- Não misture dados de diferentes ambientes

## Troubleshooting

### Erro de Conexão com Banco
```bash
# Verificar se o banco está rodando
make docker-up

# Verificar variáveis de ambiente
cat .env
```

### Erro de Migration
```bash
# Executar migrations primeiro
make migrate-up

# Depois executar seeds
make seed
```

### Limpar e Recriar Tudo
```bash
# Limpar dados seedados
make seed-clear

# Reverter migrations
make migrate-down

# Executar tudo novamente
make dev-setup
```

## Ambiente de Desenvolvimento

Para configurar um ambiente completo de desenvolvimento:

```bash
# 1. Subir containers
make docker-up

# 2. Executar migrations
make migrate-up

# 3. Executar seeds
make seed

# Ou tudo de uma vez
make dev-setup
```

Agora você tem um banco de dados completo com estrutura e dados de exemplo! 🚀
