package user

import (
	"regexp"
	"strings"
)

// PasswordValidator valida senhas seguindo boas práticas de segurança
type PasswordValidator struct {
	commonPasswords map[string]bool
	minLength       int
	maxLength       int
}

// NewPasswordValidator cria uma nova instância do validador de senhas
func NewPasswordValidator() *PasswordValidator {
	// Lista de senhas comuns que devem ser rejeitadas
	commonPasswords := map[string]bool{
		"password":    true,
		"123456":      true,
		"123456789":   true,
		"qwerty":      true,
		"abc123":      true,
		"password123": true,
		"admin":       true,
		"letmein":     true,
		"welcome":     true,
		"monkey":      true,
		"1234567890":  true,
		"password1":   true,
		"qwerty123":   true,
		"dragon":      true,
		"master":      true,
		"12345678":    true,
		"12345":       true,
		"1234567":     true,
		"123123":      true,
		"111111":      true,
		"123454":      true,
		"123453":      true,
		"123452":      true,
		"123451":      true,
		"123450":      true,
	}

	return &PasswordValidator{
		commonPasswords: commonPasswords,
		minLength:       8,
		maxLength:       128,
	}
}

// IsCommon verifica se a senha é muito comum
func (v *PasswordValidator) IsCommon(password string) bool {
	return v.commonPasswords[strings.ToLower(password)]
}

// IsValidLength verifica se a senha tem tamanho válido
func (v *PasswordValidator) IsValidLength(password string) bool {
	length := len(password)
	return length >= v.minLength && length <= v.maxLength
}

// HasRequiredChars verifica se a senha tem caracteres obrigatórios
func (v *PasswordValidator) HasRequiredChars(password string) bool {
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`).MatchString(password)

	// Pelo menos 3 dos 4 tipos de caracteres
	types := 0
	if hasLower {
		types++
	}
	if hasUpper {
		types++
	}
	if hasDigit {
		types++
	}
	if hasSpecial {
		types++
	}

	return types >= 3
}

// CalculateSimilarity calcula similaridade entre duas senhas usando distância de Levenshtein
func (v *PasswordValidator) CalculateSimilarity(s1, s2 string) float64 {
	distance := v.levenshteinDistance(s1, s2)
	maxLen := v.max(len(s1), len(s2))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(distance)/float64(maxLen)
}

// IsTooSimilar verifica se as senhas são muito similares
func (v *PasswordValidator) IsTooSimilar(oldPassword, newPassword string, threshold float64) bool {
	similarity := v.CalculateSimilarity(oldPassword, newPassword)
	return similarity > threshold
}

// Validate valida uma senha completamente
func (v *PasswordValidator) Validate(password string) error {
	if !v.IsValidLength(password) {
		return NewValidationError("PASSWORD_LENGTH", "password must be between 8 and 128 characters")
	}

	if v.IsCommon(password) {
		return NewValidationError("PASSWORD_COMMON", "password is too common, please choose a stronger password")
	}

	if !v.HasRequiredChars(password) {
		return NewValidationError("PASSWORD_CHARS", "password must contain at least 3 of: lowercase, uppercase, digits, special characters")
	}

	return nil
}

// ValidateChange valida mudança de senha
func (v *PasswordValidator) ValidateChange(oldPassword, newPassword string) error {
	// Validar nova senha
	if err := v.Validate(newPassword); err != nil {
		return err
	}

	// Verificar se não é muito similar à senha atual
	if v.IsTooSimilar(oldPassword, newPassword, 0.8) {
		return NewValidationError("PASSWORD_SIMILAR", "new password is too similar to current password")
	}

	// Verificar se não é a mesma senha
	if oldPassword == newPassword {
		return NewValidationError("PASSWORD_SAME", "new password must be different from current password")
	}

	return nil
}

// levenshteinDistance calcula a distância de Levenshtein entre duas strings
func (v *PasswordValidator) levenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	rows := len(r1) + 1
	cols := len(r2) + 1

	d := make([][]int, rows)
	for i := range d {
		d[i] = make([]int, cols)
	}

	for i := 1; i < rows; i++ {
		d[i][0] = i
	}
	for j := 1; j < cols; j++ {
		d[0][j] = j
	}

	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			d[i][j] = v.min(
				v.min(d[i-1][j]+1, d[i][j-1]+1), // min of deletion and insertion
				d[i-1][j-1]+cost,                // substitution
			)
		}
	}

	return d[rows-1][cols-1]
}

// min retorna o menor valor entre dois inteiros
func (v *PasswordValidator) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max retorna o maior valor entre dois inteiros
func (v *PasswordValidator) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NameValidator valida nomes de usuários
type NameValidator struct {
	dangerousCharsRegex *regexp.Regexp
	minLength           int
	maxLength           int
}

// NewNameValidator cria uma nova instância do validador de nomes
func NewNameValidator() *NameValidator {
	// Regex para detectar caracteres perigosos (XSS, injection, etc)
	regex := regexp.MustCompile(`[<>&"'\\\/=+()[\]{}|` + "`" + `~!@#$%^*?:;,.]`)

	return &NameValidator{
		dangerousCharsRegex: regex,
		minLength:           2,
		maxLength:           100,
	}
}

// ContainsDangerousChars verifica se o nome contém caracteres perigosos
func (v *NameValidator) ContainsDangerousChars(name string) bool {
	return v.dangerousCharsRegex.MatchString(name)
}

// IsValidLength verifica se o nome tem tamanho válido
func (v *NameValidator) IsValidLength(name string) bool {
	length := len(strings.TrimSpace(name))
	return length >= v.minLength && length <= v.maxLength
}

// IsValid verifica se o nome é válido completamente
func (v *NameValidator) IsValid(name string) bool {
	trimmed := strings.TrimSpace(name)

	if !v.IsValidLength(trimmed) {
		return false
	}

	if v.ContainsDangerousChars(trimmed) {
		return false
	}

	return true
}

// Validate valida um nome completamente
func (v *NameValidator) Validate(name string) error {
	trimmed := strings.TrimSpace(name)

	if !v.IsValidLength(trimmed) {
		return NewValidationError("NAME_LENGTH", "name must be between 2 and 100 characters")
	}

	if v.ContainsDangerousChars(trimmed) {
		return NewValidationError("NAME_DANGEROUS_CHARS", "name contains dangerous characters")
	}

	return nil
}

// EmailValidator valida emails (usando regex básico + validação do domain layer)
type EmailValidator struct {
	emailRegex *regexp.Regexp
}

// NewEmailValidator cria uma nova instância do validador de emails
func NewEmailValidator() *EmailValidator {
	// Regex básico para validação de formato de email
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return &EmailValidator{
		emailRegex: regex,
	}
}

// IsValidFormat verifica se o email tem formato válido
func (v *EmailValidator) IsValidFormat(email string) bool {
	return v.emailRegex.MatchString(email)
}

// Validate valida um email completamente
func (v *EmailValidator) Validate(email string) error {
	if !v.IsValidFormat(email) {
		return NewValidationError("EMAIL_FORMAT", "invalid email format")
	}

	// Validação adicional de comprimento
	if len(email) > 254 {
		return NewValidationError("EMAIL_LENGTH", "email too long")
	}

	return nil
}

// PhoneValidator valida números de telefone
type PhoneValidator struct {
	phoneRegex *regexp.Regexp
}

// NewPhoneValidator cria uma nova instância do validador de telefones
func NewPhoneValidator() *PhoneValidator {
	// Regex para telefones brasileiros (formato básico)
	regex := regexp.MustCompile(`^(\+55\s?)?(\(?\d{2}\)?\s?)?\d{4,5}-?\d{4}$`)

	return &PhoneValidator{
		phoneRegex: regex,
	}
}

// IsValidFormat verifica se o telefone tem formato válido
func (v *PhoneValidator) IsValidFormat(phone string) bool {
	return v.phoneRegex.MatchString(phone)
}

// Validate valida um telefone completamente
func (v *PhoneValidator) Validate(phone string) error {
	if phone == "" {
		return nil // Telefone é opcional
	}

	if !v.IsValidFormat(phone) {
		return NewValidationError("PHONE_FORMAT", "invalid phone format")
	}

	return nil
}

// ValidationError representa um erro de validação
type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implementa a interface error
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError cria um novo erro de validação
func NewValidationError(code, message string) *ValidationError {
	return &ValidationError{
		Code:    code,
		Message: message,
	}
}
