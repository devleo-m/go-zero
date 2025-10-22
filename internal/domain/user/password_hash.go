package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordHashConfig configuração para hash de senha usando Argon2
type PasswordHashConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultPasswordHashConfig retorna configuração padrão para hash de senha
func DefaultPasswordHashConfig() PasswordHashConfig {
	return PasswordHashConfig{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// hashPassword cria hash da senha usando Argon2
func hashPassword(password string) (string, error) {
	config := DefaultPasswordHashConfig()

	// Gerar salt aleatório
	salt := make([]byte, config.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Gerar hash usando Argon2
	hash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// Codificar em base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Retornar formato: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		config.Memory,
		config.Iterations,
		config.Parallelism,
		b64Salt,
		b64Hash,
	), nil
}

// verifyPassword verifica se a senha corresponde ao hash
func verifyPassword(password, encodedHash string) bool {
	// Parsear o hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}

	// Verificar algoritmo
	if parts[1] != "argon2id" {
		return false
	}

	// Parsear parâmetros
	var version int
	var memory, iterations uint32
	var parallelism uint8

	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false
	}

	// Decodificar salt e hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	// Gerar hash da senha fornecida
	otherHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	// Comparar usando constant time comparison
	return subtle.ConstantTimeCompare(hash, otherHash) == 1
}
