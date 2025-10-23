package user

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Exemplo de uso da Infrastructure Layer
// Este arquivo demonstra como usar o repository de usuários

func ExampleUsage() {
	// Este é um exemplo de como usar a infrastructure layer
	// Em produção, você teria as dependências reais injetadas

	// 1. Setup do banco (exemplo)
	var db *gorm.DB
	var logger Logger

	// 2. Criar repository
	_ = NewRepository(db, logger)

	// 3. Exemplos de uso
	_ = context.Background()

	// Criar usuário
	// user := &user.User{...}
	// err := repo.Create(ctx, user)

	// Buscar por ID
	// foundUser, err := repo.FindByID(ctx, userID)

	// Buscar por email
	// foundUser, err := repo.FindByEmail(ctx, "user@example.com")

	// Buscar com filtros
	// filter := shared.NewQueryBuilder().
	//     WhereEqual("status", "active").
	//     OrderByDesc("created_at").
	//     Page(1).
	//     PageSize(20).
	//     Build()
	// users, err := repo.FindMany(ctx, filter)

	// Queries específicas
	// result, err := repo.FindActiveUsers(ctx, 1, 20)
	// result, err := repo.FindUsersByRole(ctx, "admin", 1, 20)
	// result, err := repo.SearchUsers(ctx, "João", 1, 20)

	// Transações
	// err := repo.WithTransaction(ctx, func(txRepo *Repository) error {
	//     // Operações dentro da transação
	//     return nil
	// })

	// Estatísticas
	// stats, err := repo.GetStats(ctx)

	// Manutenção
	// err := repo.CleanupExpiredTokens(ctx)

	log.Println("Exemplo de uso da infrastructure layer")
}

// MockLogger implementa a interface Logger para testes
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...interface{}) {
	log.Printf("[WARN] %s %v", msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

// Exemplo de teste unitário
func ExampleTest() {
	// Setup de teste
	_ = &MockLogger{}

	// Em um teste real, você usaria um banco de teste
	// db := setupTestDB()
	// repo := NewRepository(db, mockLogger)

	// Testes de CRUD
	// TestCreateUser(t, repo)
	// TestFindUserByID(t, repo)
	// TestFindUserByEmail(t, repo)
	// TestUpdateUser(t, repo)
	// TestDeleteUser(t, repo)

	// Testes de queries
	// TestFindActiveUsers(t, repo)
	// TestSearchUsers(t, repo)
	// TestGetStats(t, repo)

	log.Println("Exemplo de teste da infrastructure layer")
}

// Exemplo de migração do banco
func ExampleMigration() {
	// SQL para criar as tabelas
	sql := `
-- Tabela principal de usuários
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    role VARCHAR(50) NOT NULL DEFAULT 'user'
);

-- Índices
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Tabela de perfil
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    email_verified_at TIMESTAMP,
    last_login_at TIMESTAMP,
    login_count INTEGER DEFAULT 0 NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    avatar_url VARCHAR(500),
    bio TEXT,
    location VARCHAR(255),
    website VARCHAR(500)
);

-- Tabela de dados de autenticação
CREATE TABLE user_auth_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    password_reset_token VARCHAR(255) UNIQUE,
    password_reset_expires TIMESTAMP,
    activation_token VARCHAR(255) UNIQUE,
    activation_expires TIMESTAMP,
    two_factor_secret VARCHAR(255),
    two_factor_enabled BOOLEAN DEFAULT false NOT NULL,
    two_factor_backup_codes TEXT,
    refresh_token VARCHAR(500) UNIQUE,
    refresh_expires TIMESTAMP,
    refresh_token_hash VARCHAR(255),
    last_password_change TIMESTAMP,
    failed_login_attempts INTEGER DEFAULT 0 NOT NULL,
    locked_until TIMESTAMP
);

-- Tabela de preferências
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    timezone VARCHAR(50) DEFAULT 'UTC',
    language VARCHAR(10) DEFAULT 'pt-BR',
    currency VARCHAR(3) DEFAULT 'BRL',
    email_notifications BOOLEAN DEFAULT true NOT NULL,
    sms_notifications BOOLEAN DEFAULT false NOT NULL,
    push_notifications BOOLEAN DEFAULT true NOT NULL,
    marketing_emails BOOLEAN DEFAULT false NOT NULL,
    security_alerts BOOLEAN DEFAULT true NOT NULL,
    appointment_reminders BOOLEAN DEFAULT true NOT NULL,
    newsletter_subscription BOOLEAN DEFAULT false NOT NULL,
    profile_visibility VARCHAR(20) DEFAULT 'private' NOT NULL,
    show_email BOOLEAN DEFAULT false NOT NULL,
    show_phone BOOLEAN DEFAULT false NOT NULL,
    show_last_login BOOLEAN DEFAULT false NOT NULL,
    theme VARCHAR(20) DEFAULT 'light' NOT NULL,
    date_format VARCHAR(20) DEFAULT 'DD/MM/YYYY' NOT NULL,
    time_format VARCHAR(10) DEFAULT '24h' NOT NULL,
    items_per_page INTEGER DEFAULT 20 NOT NULL,
    auto_save_drafts BOOLEAN DEFAULT true NOT NULL,
    show_tutorials BOOLEAN DEFAULT true NOT NULL,
    high_contrast BOOLEAN DEFAULT false NOT NULL,
    large_text BOOLEAN DEFAULT false NOT NULL,
    screen_reader BOOLEAN DEFAULT false NOT NULL,
    keyboard_nav BOOLEAN DEFAULT false NOT NULL,
    reduced_motion BOOLEAN DEFAULT false NOT NULL
);
`

	log.Println("SQL de migração:")
	fmt.Println(sql)
}

// Exemplo de configuração do GORM
func ExampleGORMConfig() {
	// Configuração recomendada do GORM
	config := &gorm.Config{
		// NamingStrategy: schema.NamingStrategy{
		//     TablePrefix:   "app_",
		//     SingularTable: false,
		// },
		// Logger: logger.Default.LogMode(logger.Info),
		// NowFunc: func() time.Time {
		//     return time.Now().UTC()
		// },
		// DisableForeignKeyConstraintWhenMigrating: true,
		// PrepareStmt: true,
	}

	log.Printf("Configuração do GORM: %+v", config)
}

// Exemplo de métricas
func ExampleMetrics() {
	// Métricas recomendadas para monitoramento
	metrics := map[string]string{
		"user_operations_total":       "Total de operações por tipo",
		"user_query_duration_seconds": "Tempo de queries",
		"user_cache_hits_total":       "Cache hits (se implementado)",
		"user_errors_total":           "Erros por tipo",
		"user_active_count":           "Usuários ativos",
		"user_registrations_total":    "Total de registros",
		"user_logins_total":           "Total de logins",
		"user_password_resets_total":  "Total de resets de senha",
	}

	log.Println("Métricas recomendadas:")
	for metric, description := range metrics {
		log.Printf("- %s: %s", metric, description)
	}
}
