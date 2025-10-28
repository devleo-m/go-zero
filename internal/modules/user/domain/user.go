package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User representa um usuário no domínio.
type User struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Phone     *string    `json:"phone,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	ID        uuid.UUID  `json:"id"`
}

// NewUser cria um novo usuário.
func NewUser(name, email, password string) (*User, error) {
	// Validações básicas
	if name == "" || len(name) < 2 {
		return nil, ErrInvalidName
	}

	// Validação de email
	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, ErrInvalidEmail
	}

	if password == "" || len(password) < 8 {
		return nil, ErrInvalidPassword
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrPasswordHash
	}

	now := time.Now()

	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      "user",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ValidatePassword verifica se a senha está correta.
func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// UpdatePassword atualiza a senha do usuário.
func (u *User) UpdatePassword(newPassword string) error {
	if newPassword == "" || len(newPassword) < 8 {
		return ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return ErrPasswordHash
	}

	u.Password = string(hashedPassword)
	u.UpdatedAt = time.Now()

	return nil
}

// UpdateProfile atualiza informações do perfil.
func (u *User) UpdateProfile(name string, phone *string) error {
	if name == "" || len(name) < 2 {
		return ErrInvalidName
	}

	u.Name = name
	u.Phone = phone
	u.UpdatedAt = time.Now()

	return nil
}

// SoftDelete marca o usuário como deletado.
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}

// IsDeleted verifica se o usuário foi deletado.
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}
