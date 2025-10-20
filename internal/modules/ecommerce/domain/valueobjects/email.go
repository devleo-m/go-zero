package valueobjects

import (
	"regexp"
	"strings"

	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/errors"
)

// Email representa um email válido
type Email struct {
	value string
}

// NewEmail cria um novo Email value object
func NewEmail(email string) (*Email, error) {
	if email == "" {
		return nil, errors.ErrInvalidEmail
	}

	// Normalizar email (lowercase, trim)
	email = strings.ToLower(strings.TrimSpace(email))

	// Validar formato
	if !isValidEmail(email) {
		return nil, errors.ErrInvalidEmail
	}

	return &Email{value: email}, nil
}

// String retorna o valor do email
func (e Email) String() string {
	return e.value
}

// Value retorna o valor do email (alias para String)
func (e Email) Value() string {
	return e.value
}

// Equals verifica se dois emails são iguais
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsValid verifica se o email é válido
func (e Email) IsValid() bool {
	return isValidEmail(e.value)
}

// GetDomain retorna o domínio do email
func (e Email) GetDomain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// GetLocalPart retorna a parte local do email (antes do @)
func (e Email) GetLocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// isValidEmail valida o formato do email usando regex
func isValidEmail(email string) bool {
	// Regex para validar email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

	if !emailRegex.MatchString(email) {
		return false
	}

	// Validações adicionais
	if len(email) > 254 { // RFC 5321
		return false
	}

	// Verificar se não começa ou termina com ponto
	if strings.HasPrefix(email, ".") || strings.HasSuffix(email, ".") {
		return false
	}

	// Verificar se não tem pontos consecutivos
	if strings.Contains(email, "..") {
		return false
	}

	return true
}

// MustNewEmail cria um email sem retornar erro (para casos onde sabemos que é válido)
func MustNewEmail(email string) Email {
	e, err := NewEmail(email)
	if err != nil {
		panic(err)
	}
	return *e
}
