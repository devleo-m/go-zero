package seeds

import (
	"log"

	"github.com/devleo-m/go-zero/internal/domain/user"
	postgresUser "github.com/devleo-m/go-zero/internal/infrastructure/persistence/postgres/user"
	"gorm.io/gorm"
)

// UserSeeder handles seeding user data
type UserSeeder struct {
	db *gorm.DB
}

// NewUserSeeder creates a new user seeder
func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{
		db: db,
	}
}

// SeedUsers seeds the database with initial user data
func (s *UserSeeder) SeedUsers() error {
	log.Println("🌱 Starting user seeding...")

	// Define users to seed
	usersToSeed := []struct {
		name     string
		email    string
		password string
		role     user.Role
		status   user.Status
		phone    string
	}{
		{
			name:     "Administrador do Sistema",
			email:    "admin@go-zero.com",
			password: "Admin123!@#",
			role:     user.RoleAdmin,
			status:   user.StatusActive,
			phone:    "+5511999999999",
		},
		{
			name:     "Gerente de Projetos",
			email:    "manager@go-zero.com",
			password: "Manager123!@#",
			role:     user.RoleManager,
			status:   user.StatusActive,
			phone:    "+5511888888888",
		},
		{
			name:     "Usuário Teste",
			email:    "user@go-zero.com",
			password: "User123!@#",
			role:     user.RoleUser,
			status:   user.StatusActive,
			phone:    "+5511777777777",
		},
		{
			name:     "Usuário Pendente",
			email:    "pending@go-zero.com",
			password: "Pending123!@#",
			role:     user.RoleUser,
			status:   user.StatusPending,
			phone:    "+5511666666666",
		},
		{
			name:     "Usuário Suspenso",
			email:    "suspended@go-zero.com",
			password: "Suspended123!@#",
			role:     user.RoleUser,
			status:   user.StatusSuspended,
			phone:    "+5511555555555",
		},
		{
			name:     "Convidado",
			email:    "guest@go-zero.com",
			password: "Guest123!@#",
			role:     user.RoleGuest,
			status:   user.StatusActive,
		},
	}

	// Seed each user
	for _, userData := range usersToSeed {
		if err := s.seedUser(userData); err != nil {
			log.Printf("❌ Error seeding user %s: %v", userData.email, err)
			return err
		}
	}

	log.Println("✅ User seeding completed successfully!")
	return nil
}

// seedUser seeds a single user
func (s *UserSeeder) seedUser(userData struct {
	name     string
	email    string
	password string
	role     user.Role
	status   user.Status
	phone    string
}) error {
	// Check if user already exists
	var existingUser postgresUser.UserModel
	if err := s.db.Where("email = ?", userData.email).First(&existingUser).Error; err == nil {
		log.Printf("⚠️  User %s already exists, skipping...", userData.email)
		return nil
	}

	// Create domain entities
	emailVO, err := user.NewEmail(userData.email)
	if err != nil {
		return err
	}

	passwordVO, err := user.NewPassword(userData.password)
	if err != nil {
		return err
	}

	// Create user domain entity
	userEntity, err := user.NewUser(userData.name, emailVO, passwordVO, userData.role)
	if err != nil {
		return err
	}

	// Set status manually (since NewUser always sets to pending)
	userEntity.Status = userData.status

	// Set phone if provided
	if userData.phone != "" {
		phoneVO, err := user.NewPhone(userData.phone)
		if err != nil {
			return err
		}
		userEntity.Phone = &phoneVO
	}

	// Convert to GORM model
	userModel := postgresUser.ToModel(userEntity)

	// Create in database
	if err := s.db.Create(&userModel).Error; err != nil {
		return err
	}

	log.Printf("✅ Created user: %s (%s)", userData.name, userData.email)
	return nil
}

// ClearUsers removes all seeded users (for testing)
func (s *UserSeeder) ClearUsers() error {
	log.Println("🧹 Clearing seeded users...")

	// Delete users with seeded emails
	seededEmails := []string{
		"admin@go-zero.com",
		"manager@go-zero.com",
		"user@go-zero.com",
		"pending@go-zero.com",
		"suspended@go-zero.com",
		"guest@go-zero.com",
	}

	result := s.db.Where("email IN ?", seededEmails).Delete(&postgresUser.UserModel{})
	if result.Error != nil {
		return result.Error
	}

	log.Printf("✅ Cleared %d seeded users", result.RowsAffected)
	return nil
}
