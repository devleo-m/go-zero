package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devleo-m/go-zero/internal/infrastructure"
	userRepo "github.com/devleo-m/go-zero/internal/modules/user/infrastructure/postgres"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Conectar ao banco de dados
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=go_zero port=5432 sslmode=disable"
	}

	db, err := infrastructure.NewDatabase(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Executar seeds de lookup tables
	log.Println("🌱 Executando seeds de lookup tables...")
	if err := executeLookupSeeds(db.DB); err != nil {
		log.Fatal("Failed to execute lookup seeds:", err)
	}

	// Executar migração de dados de usuários
	log.Println("🔄 Migrando dados de usuários...")
	if err := executeUserMigration(db.DB); err != nil {
		log.Fatal("Failed to migrate user data:", err)
	}

	// Configurar repositório
	userRepository := userRepo.NewRepository(db.DB)

	// Criar usuário de exemplo
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	log.Println("✅ Todos os seeds executados com sucesso!")
}

// executeLookupSeeds executa os seeds das tabelas de lookup
func executeLookupSeeds(db *sql.DB) error {
	seedDir := "database/seeds"

	// Listar arquivos SQL na pasta seeds
	files, err := ioutil.ReadDir(seedDir)
	if err != nil {
		return fmt.Errorf("failed to read seeds directory: %w", err)
	}

	// Filtrar apenas arquivos .sql e ordenar
	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") && !strings.Contains(file.Name(), "migrate_users_data") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Executar cada arquivo SQL
	for _, fileName := range sqlFiles {
		filePath := filepath.Join(seedDir, fileName)
		log.Printf("📄 Executando seed: %s", fileName)

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %w", fileName, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute seed %s: %w", fileName, err)
		}
	}

	return nil
}

// executeUserMigration executa a migração de dados de usuários
func executeUserMigration(db *sql.DB) error {
	seedFile := "database/seeds/06_migrate_users_data.sql"

	log.Printf("📄 Executando migração: %s", seedFile)

	content, err := ioutil.ReadFile(seedFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute user migration: %w", err)
	}

	return nil
}
