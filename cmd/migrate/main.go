package main

import (
	"flag"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/go_zero?sslmode=disable"
	migrationsPath     = "file://database/migrations"
	directionUp        = "up"
	directionDown      = "down"
	directionForce     = "force"
)

func main() {
	// Flags
	direction, steps := parseFlags()

	// Obter URL do banco
	databaseURL := getDatabaseURL()

	// Executar migration
	executeMigration(databaseURL, direction, steps)
}

// parseFlags parse os arguments da linha de commando.
func parseFlags() (string, int) {
	var direction string

	var steps int

	flag.StringVar(&direction, "direction", directionUp, "Migration direction: up, down, force")
	flag.IntVar(&steps, "steps", 0, "Number of steps (0 = all)")
	flag.Parse()

	return direction, steps
}

// getDatabaseURL obtém a URL do banco de dados.
func getDatabaseURL() string {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}

	return databaseURL
}

// executeMigration executa a migration.
func executeMigration(databaseURL, direction string, steps int) {
	// Criar migrator
	m, err := createMigrator(databaseURL)
	if err != nil {
		log.Fatalf("❌ Erro ao criar migrator: %v", err)
	}
	defer closeMigrator(m)

	// Executar migration baseado na direção
	err = runMigration(m, direction, steps)
	if err != nil {
		handleMigrationError(err)
		return
	}

	// Log de sucesso
	logSuccess(m, direction)
}

// createMigrator cria e retorna um migrator.
func createMigrator(databaseURL string) (*migrate.Migrate, error) {
	return migrate.New(migrationsPath, databaseURL)
}

// closeMigrator fecha o migrator.
func closeMigrator(m *migrate.Migrate) {
	if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
		log.Printf("⚠️  Erro ao fechar migrator - Source: %v, DB: %v", srcErr, dbErr)
	}
}

// runMigration executa a migration baseado na direção.
func runMigration(m *migrate.Migrate, direction string, steps int) error {
	switch direction {
	case directionUp:
		return runUp(m, steps)
	case directionDown:
		return runDown(m, steps)
	case directionForce:
		return runForce(m, steps)
	default:
		log.Fatal("❌ Direction inválida: use up, down ou force")
		return nil
	}
}

// runUp executa migration para cima.
func runUp(m *migrate.Migrate, steps int) error {
	if steps > 0 {
		return m.Steps(steps)
	}

	return m.Up()
}

// runDown executa migration para baixo.
func runDown(m *migrate.Migrate, steps int) error {
	if steps > 0 {
		return m.Steps(-steps)
	}

	return m.Down()
}

// runForce força uma versão específica.
func runForce(m *migrate.Migrate, steps int) error {
	if steps == 0 {
		log.Fatal("❌ Especifique a versão com -steps")
	}

	return m.Force(steps)
}

// handleMigrationError trata erros de migration.
func handleMigrationError(err error) {
	if err != migrate.ErrNoChange {
		log.Fatalf("❌ Migration falhou: %v", err)
	}
}

// logSuccess registra sucesso da migration.
func logSuccess(m *migrate.Migrate, direction string) {
	version, dirty, err := m.Version()
	if err != nil {
		log.Printf("⚠️  Não foi possível obter versão: %v", err)
		return
	}

	log.Printf("✅ Migration executada com sucesso! Versão: %d, Dirty: %v", version, dirty)
	log.Printf("🎉 Migration %s executada! Versão: %d", direction, version)
}
