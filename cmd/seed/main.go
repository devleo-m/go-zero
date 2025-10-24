package main

import (
	"log"
	"os"

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

	// Configurar repositório
	userRepository := userRepo.NewRepository(db.DB)

	// Criar usuário de exemplo
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	log.Println("✅ Seed executado com sucesso!")
}
