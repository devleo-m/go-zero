package services

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"
)

// PasswordService fornece serviços relacionados a senhas
type PasswordService struct{}

// NewPasswordService cria um novo serviço de senhas
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// GenerateRandomPassword gera uma senha aleatória segura
func (s *PasswordService) GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}
	if length > 128 {
		length = 128
	}

	// Caracteres permitidos
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	allChars := lowercase + uppercase + numbers + symbols

	// Garantir que a senha tenha pelo menos um de cada tipo
	password := make([]byte, length)

	// Primeiro caractere: minúscula
	password[0] = lowercase[randomInt(len(lowercase))]

	// Segundo caractere: maiúscula
	password[1] = uppercase[randomInt(len(uppercase))]

	// Terceiro caractere: número
	password[2] = numbers[randomInt(len(numbers))]

	// Quarto caractere: símbolo
	password[3] = symbols[randomInt(len(symbols))]

	// Resto: aleatório
	for i := 4; i < length; i++ {
		password[i] = allChars[randomInt(len(allChars))]
	}

	// Embaralhar a senha
	for i := len(password) - 1; i > 0; i-- {
		j := randomInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// ValidatePasswordStrength valida a força de uma senha
func (s *PasswordService) ValidatePasswordStrength(password string) (bool, []string) {
	// Usar o value object para validação (evita duplicação)
	_, err := valueobjects.NewPassword(password)
	if err != nil {
		return false, []string{err.Error()}
	}

	// Validações adicionais específicas do serviço
	var errors []string

	// Verificar sequências comuns
	if s.hasCommonSequences(password) {
		errors = append(errors, "senha não deve conter sequências comuns (123, abc, etc.)")
	}

	// Verificar repetições
	if s.hasRepeatedCharacters(password) {
		errors = append(errors, "senha não deve conter caracteres repetidos em sequência")
	}

	return len(errors) == 0, errors
}

// GeneratePasswordHash gera hash de uma senha
func (s *PasswordService) GeneratePasswordHash(password string) (string, error) {
	passwordVO, err := valueobjects.NewPassword(password)
	if err != nil {
		return "", err
	}

	return passwordVO.Hash(), nil
}

// VerifyPassword verifica se uma senha corresponde ao hash
func (s *PasswordService) VerifyPassword(password, hash string) bool {
	passwordVO, err := valueobjects.NewPasswordFromHash(hash)
	if err != nil {
		return false
	}

	return passwordVO.Verify(password)
}

// hasCommonSequences verifica se a senha tem sequências comuns
func (s *PasswordService) hasCommonSequences(password string) bool {
	commonSequences := []string{
		"123", "234", "345", "456", "567", "678", "789",
		"abc", "bcd", "cde", "def", "efg", "fgh", "ghi",
		"qwe", "wer", "ert", "rty", "tyu", "yui", "uio",
		"asd", "sdf", "dfg", "fgh", "ghj", "hjk", "jkl",
		"zxc", "xcv", "cvb", "vbn", "bnm",
		"password", "123456", "qwerty", "admin", "user",
	}

	passwordLower := strings.ToLower(password)
	for _, seq := range commonSequences {
		if strings.Contains(passwordLower, seq) {
			return true
		}
	}

	return false
}

// hasRepeatedCharacters verifica se a senha tem caracteres repetidos em sequência
func (s *PasswordService) hasRepeatedCharacters(password string) bool {
	if len(password) < 3 {
		return false
	}

	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}

	return false
}

// randomInt gera um número aleatório entre 0 e max-1
func randomInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}
