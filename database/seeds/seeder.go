package seeds

import (
	"log"

	"gorm.io/gorm"
)

// SeederManager manages all seeders
type SeederManager struct {
	db *gorm.DB
}

// NewSeederManager creates a new seeder manager
func NewSeederManager(db *gorm.DB) *SeederManager {
	return &SeederManager{
		db: db,
	}
}

// RunAll executes all seeders
func (sm *SeederManager) RunAll() error {
	log.Println("🌱 Starting database seeding...")
	log.Println("=====================================")

	// Initialize seeders
	userSeeder := NewUserSeeder(sm.db)

	// Run seeders in order
	if err := userSeeder.SeedUsers(); err != nil {
		log.Printf("❌ User seeding failed: %v", err)
		return err
	}

	log.Println("=====================================")
	log.Println("✅ All seeders completed successfully!")
	return nil
}

// RunUserSeeder executes only the user seeder
func (sm *SeederManager) RunUserSeeder() error {
	log.Println("🌱 Starting user seeding...")

	userSeeder := NewUserSeeder(sm.db)
	return userSeeder.SeedUsers()
}

// ClearAll clears all seeded data
func (sm *SeederManager) ClearAll() error {
	log.Println("🧹 Clearing all seeded data...")
	log.Println("=====================================")

	// Initialize seeders
	userSeeder := NewUserSeeder(sm.db)

	// Clear seeders in reverse order
	if err := userSeeder.ClearUsers(); err != nil {
		log.Printf("❌ User clearing failed: %v", err)
		return err
	}

	log.Println("=====================================")
	log.Println("✅ All seeded data cleared successfully!")
	return nil
}

// ClearUserSeeder clears only user seeded data
func (sm *SeederManager) ClearUserSeeder() error {
	log.Println("🧹 Clearing user seeded data...")

	userSeeder := NewUserSeeder(sm.db)
	return userSeeder.ClearUsers()
}
