package valueobjects

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/errors"
	"golang.org/x/crypto/argon2"
)

// Password representa uma senha segura
type Password struct {
	hash string
}

// PasswordConfig configuração para hash de senha
type PasswordConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultPasswordConfig configuração padrão para hash de senha
var DefaultPasswordConfig = &PasswordConfig{
	Memory:      64 * 1024, // 64 MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// NewPassword cria uma nova senha com hash
func NewPassword(plainPassword string) (*Password, error) {
	if plainPassword == "" {
		return nil, errors.ErrInvalidPassword
	}

	// Validar força da senha
	if !isStrongPassword(plainPassword) {
		return nil, errors.ErrInvalidPassword
	}

	// Gerar hash
	hash, err := hashPassword(plainPassword, DefaultPasswordConfig)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	return &Password{hash: hash}, nil
}

// NewPasswordFromHash cria uma senha a partir de um hash existente
func NewPasswordFromHash(hash string) (*Password, error) {
	if hash == "" {
		return nil, errors.ErrInvalidPassword
	}

	// Validar formato do hash
	if !isValidHash(hash) {
		return nil, errors.ErrInvalidPassword
	}

	return &Password{hash: hash}, nil
}

// Hash retorna o hash da senha
func (p Password) Hash() string {
	return p.hash
}

// Verify verifica se a senha em texto plano corresponde ao hash
func (p Password) Verify(plainPassword string) bool {
	// Decodificar hash
	hash, salt, config, err := decodeHash(p.hash)
	if err != nil {
		return false
	}

	// Gerar hash da senha fornecida
	otherHash := argon2.IDKey([]byte(plainPassword), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// Comparar hashes de forma segura
	return subtle.ConstantTimeCompare(hash, otherHash) == 1
}

// Equals verifica se duas senhas são iguais (comparando hashes)
func (p Password) Equals(other Password) bool {
	return p.hash == other.hash
}

// IsValid verifica se a senha é válida
func (p Password) IsValid() bool {
	return isValidHash(p.hash)
}

// hashPassword gera hash da senha usando Argon2
func hashPassword(password string, config *PasswordConfig) (string, error) {
	// Gerar salt aleatório
	salt := make([]byte, config.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Gerar hash
	hash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// Codificar em base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Retornar hash formatado
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Iterations, config.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// decodeHash decodifica um hash para extrair componentes
func decodeHash(encodedHash string) (hash, salt []byte, config *PasswordConfig, err error) {
	// Verificar formato
	if !strings.HasPrefix(encodedHash, "$argon2id$") {
		return nil, nil, nil, fmt.Errorf("formato de hash inválido")
	}

	// Parse do hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, fmt.Errorf("formato de hash inválido")
	}

	// Decodificar salt e hash
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}

	// Parse da configuração
	config = &PasswordConfig{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &config.Memory, &config.Iterations, &config.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	config.SaltLength = uint32(len(salt))
	config.KeyLength = uint32(len(hash))

	return hash, salt, config, nil
}

// isStrongPassword valida a força da senha
func isStrongPassword(password string) bool {
	// Mínimo 8 caracteres
	if len(password) < 8 {
		return false
	}

	// Máximo 128 caracteres
	if len(password) > 128 {
		return false
	}

	// Deve conter pelo menos uma letra minúscula
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return false
	}

	// Deve conter pelo menos uma letra maiúscula
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return false
	}

	// Deve conter pelo menos um número
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return false
	}

	// Deve conter pelo menos um caractere especial
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return false
	}

	return true
}

// isValidHash valida se o hash está no formato correto
func isValidHash(hash string) bool {
	return strings.HasPrefix(hash, "$argon2id$") && len(hash) > 50
}

// MustNewPassword cria uma senha sem retornar erro (para casos onde sabemos que é válida)
func MustNewPassword(plainPassword string) Password {
	p, err := NewPassword(plainPassword)
	if err != nil {
		panic(err)
	}
	return *p
}
