package user

import (
	"regexp"
	"strings"

	"github.com/devleo-m/go-zero/internal/domain/shared"
)

// Email é um Value Object que representa um endereço de email
// É imutável e sempre válido quando criado
type Email struct {
	value string
}

// NewEmail cria um novo Email validando o formato
// Retorna erro se o email for inválido
func NewEmail(value string) (Email, error) {
	email := Email{
		value: strings.ToLower(strings.TrimSpace(value)),
	}

	if err := email.validate(); err != nil {
		return Email{}, err
	}

	return email, nil
}

// String retorna a representação string do email
func (e Email) String() string {
	return e.value
}

// validate verifica se o email tem formato válido
func (e Email) validate() error {
	if e.value == "" {
		return shared.NewDomainError("INVALID_EMAIL", "email cannot be empty")
	}

	// Regex para validar formato de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(e.value) {
		return shared.NewDomainError("INVALID_EMAIL", "invalid email format")
	}

	// Validações adicionais de negócio
	if len(e.value) > 254 {
		return shared.NewDomainError("INVALID_EMAIL", "email too long")
	}

	return nil
}

// Equals compara dois emails
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Domain retorna o domínio do email (parte após @)
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart retorna a parte local do email (antes do @)
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// Password é um Value Object que representa uma senha
// Contém a senha hasheada e métodos para validação
type Password struct {
	hashedValue string
}

// NewPassword cria uma nova senha validando os critérios
func NewPassword(plainPassword string) (Password, error) {
	if err := validatePasswordStrength(plainPassword); err != nil {
		return Password{}, err
	}

	hashed, err := hashPassword(plainPassword)
	if err != nil {
		return Password{}, shared.NewDomainError("PASSWORD_HASH_ERROR", "failed to hash password")
	}

	return Password{
		hashedValue: hashed,
	}, nil
}

// NewPasswordFromHash cria um Password a partir de um hash já existente
// Usado quando carregamos do banco de dados
func NewPasswordFromHash(hashedPassword string) Password {
	return Password{
		hashedValue: hashedPassword,
	}
}

// String retorna o hash da senha (para persistência)
func (p Password) String() string {
	return p.hashedValue
}

// Verify verifica se a senha fornecida corresponde ao hash
func (p Password) Verify(plainPassword string) bool {
	return verifyPassword(plainPassword, p.hashedValue)
}

// validatePasswordStrength valida os critérios de força da senha
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return shared.NewDomainError("WEAK_PASSWORD", "password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return shared.NewDomainError("WEAK_PASSWORD", "password too long")
	}

	// Verificar se tem pelo menos uma letra minúscula
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return shared.NewDomainError("WEAK_PASSWORD", "password must contain at least one lowercase letter")
	}

	// Verificar se tem pelo menos uma letra maiúscula
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return shared.NewDomainError("WEAK_PASSWORD", "password must contain at least one uppercase letter")
	}

	// Verificar se tem pelo menos um número
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return shared.NewDomainError("WEAK_PASSWORD", "password must contain at least one number")
	}

	// Verificar se tem pelo menos um caractere especial
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return shared.NewDomainError("WEAK_PASSWORD", "password must contain at least one special character")
	}

	return nil
}

// Phone é um Value Object que representa um número de telefone
type Phone struct {
	value string
}

// NewPhone cria um novo Phone validando o formato
func NewPhone(value string) (Phone, error) {
	phone := Phone{
		value: strings.TrimSpace(value),
	}

	if err := phone.validate(); err != nil {
		return Phone{}, err
	}

	return phone, nil
}

// String retorna a representação string do telefone
func (p Phone) String() string {
	return p.value
}

// validate verifica se o telefone tem formato válido
func (p Phone) validate() error {
	if p.value == "" {
		return shared.NewDomainError("INVALID_PHONE", "phone cannot be empty")
	}

	// Remove caracteres não numéricos para validação
	digits := regexp.MustCompile(`\D`).ReplaceAllString(p.value, "")

	if len(digits) < 10 || len(digits) > 15 {
		return shared.NewDomainError("INVALID_PHONE", "phone must have between 10 and 15 digits")
	}

	return nil
}

// Equals compara dois telefones
func (p Phone) Equals(other Phone) bool {
	return p.value == other.value
}

// Format retorna o telefone formatado
func (p Phone) Format() string {
	digits := regexp.MustCompile(`\D`).ReplaceAllString(p.value, "")

	if len(digits) == 11 {
		// Formato: (XX) 9XXXX-XXXX
		return "(" + digits[0:2] + ") " + digits[2:3] + digits[3:7] + "-" + digits[7:11]
	} else if len(digits) == 10 {
		// Formato: (XX) XXXX-XXXX
		return "(" + digits[0:2] + ") " + digits[2:6] + "-" + digits[6:10]
	}

	return p.value
}
