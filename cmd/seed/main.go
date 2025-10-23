package main

import (
	"flag"
	"log"
	"os"

	"github.com/devleo-m/go-zero/database/seeds"
	"github.com/devleo-m/go-zero/internal/infrastructure/persistence/postgres"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  Warning: .env file not found: %v", err)
	}

	// Parse command line flags
	var (
		action = flag.String("action", "seed", "Action to perform: seed, clear, users, clear-users")
		help   = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Initialize database connection
	db, err := postgres.NewConnection()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Initialize seeder manager
	seederManager := seeds.NewSeederManager(db)

	// Execute action based on flag
	switch *action {
	case "seed":
		log.Println("🚀 Running all seeders...")
		if err := seederManager.RunAll(); err != nil {
			log.Fatalf("❌ Seeding failed: %v", err)
		}
	case "clear":
		log.Println("🧹 Clearing all seeded data...")
		if err := seederManager.ClearAll(); err != nil {
			log.Fatalf("❌ Clearing failed: %v", err)
		}
	case "users":
		log.Println("👥 Running user seeder...")
		if err := seederManager.RunUserSeeder(); err != nil {
			log.Fatalf("❌ User seeding failed: %v", err)
		}
	case "clear-users":
		log.Println("🧹 Clearing user seeded data...")
		if err := seederManager.ClearUserSeeder(); err != nil {
			log.Fatalf("❌ User clearing failed: %v", err)
		}
	default:
		log.Printf("❌ Unknown action: %s", *action)
		showHelp()
		os.Exit(1)
	}

	log.Println("🎉 Operation completed successfully!")
}

func showHelp() {
	log.Println("🌱 GO ZERO - Database Seeder")
	log.Println("==============================")
	log.Println("")
	log.Println("Usage:")
	log.Println("  go run cmd/seed/main.go [flags]")
	log.Println("")
	log.Println("Flags:")
	log.Println("  -action string")
	log.Println("        Action to perform (default 'seed')")
	log.Println("        Options:")
	log.Println("          seed        - Run all seeders")
	log.Println("          clear       - Clear all seeded data")
	log.Println("          users       - Run only user seeder")
	log.Println("          clear-users - Clear only user seeded data")
	log.Println("  -help")
	log.Println("        Show this help message")
	log.Println("")
	log.Println("Examples:")
	log.Println("  go run cmd/seed/main.go                    # Run all seeders")
	log.Println("  go run cmd/seed/main.go -action=users       # Run only user seeder")
	log.Println("  go run cmd/seed/main.go -action=clear       # Clear all seeded data")
	log.Println("  go run cmd/seed/main.go -action=clear-users # Clear user seeded data")
	log.Println("")
	log.Println("Environment Variables:")
	log.Println("  DB_HOST     - Database host (default: localhost)")
	log.Println("  DB_PORT     - Database port (default: 5432)")
	log.Println("  DB_USER     - Database user (default: postgres)")
	log.Println("  DB_PASSWORD - Database password (default: postgres)")
	log.Println("  DB_NAME     - Database name (default: go_zero)")
	log.Println("  DB_SSLMODE  - SSL mode (default: disable)")
}
