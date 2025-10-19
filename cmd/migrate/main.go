package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/infra/config"
	"github.com/devleo-m/go-zero/internal/infra/logger"
)

func main() {
	// Carregar .env
	_ = godotenv.Load()

	// Flags
	var direction string
	var steps int
	flag.StringVar(&direction, "direction", "up", "Migration direction: up, down, force")
	flag.IntVar(&steps, "steps", 0, "Number of steps (0 = all)")
	flag.Parse()

	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Erro ao carregar configurações: %v", err)
	}

	// Inicializar logger
	if err := logger.InitLogger(logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	}); err != nil {
		log.Fatalf("❌ Erro ao inicializar logger: %v", err)
	}
	defer logger.Sync()

	// Configurar migrator
	migrationsPath := "file://internal/infra/database/migrations"
	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	logger.Info("🚀 Iniciando migrations",
		zap.String("direction", direction),
		zap.Int("steps", steps),
		zap.String("database", cfg.Database.Name),
	)

	// Criar migrator
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		logger.Fatal("❌ Erro ao criar migrator", err)
	}
	defer m.Close()

	// Executar migration
	switch direction {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	case "down":
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
	case "force":
		if steps == 0 {
			logger.Fatal("❌ Especifique a versão com -steps", nil)
		}
		err = m.Force(steps)
	default:
		logger.Fatal("❌ Direction inválida: use up, down ou force", nil)
	}

	// Verificar resultado
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal("❌ Migration falhou", err)
	}

	// Log de sucesso
	version, dirty, _ := m.Version()
	logger.Info("✅ Migration executada com sucesso!",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty),
		zap.String("direction", direction),
	)

	fmt.Printf("🎉 Migration %s executada! Versão: %d\n", direction, version)
}
