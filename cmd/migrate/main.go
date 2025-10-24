package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Flags
	var direction string
	var steps int
	flag.StringVar(&direction, "direction", "up", "Migration direction: up, down, force")
	flag.IntVar(&steps, "steps", 0, "Number of steps (0 = all)")
	flag.Parse()

	// Obter URL do banco
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/go_zero?sslmode=disable"
	}

	// Configurar migrator
	migrationsPath := "file://database/migrations"

	log.Printf("🚀 Iniciando migration: %s", direction)

	// Criar migrator
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		log.Fatalf("❌ Erro ao criar migrator: %v", err)
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
			log.Fatal("❌ Especifique a versão com -steps")
		}
		err = m.Force(steps)
	default:
		log.Fatal("❌ Direction inválida: use up, down ou force")
	}

	// Verificar resultado
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Migration falhou: %v", err)
	}

	// Log de sucesso
	version, dirty, _ := m.Version()
	log.Printf("✅ Migration executada com sucesso! Versão: %d, Dirty: %v", version, dirty)
	fmt.Printf("🎉 Migration %s executada! Versão: %d\n", direction, version)
}
