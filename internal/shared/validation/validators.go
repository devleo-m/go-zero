package validation

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	uuidRegex  = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// ValidationError representa um erro de validação
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// ValidateEmail valida um endereço de email
func ValidateEmail(email string) error {
	if email == "" {
		return ValidationError{Field: "email", Message: "Email is required"}
	}

	if !emailRegex.MatchString(email) {
		return ValidationError{Field: "email", Message: "Invalid email format"}
	}

	return nil
}

// ValidatePassword valida uma senha
func ValidatePassword(password string) error {
	if password == "" {
		return ValidationError{Field: "password", Message: "Password is required"}
	}

	if len(password) < 8 {
		return ValidationError{Field: "password", Message: "Password must be at least 8 characters long"}
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ValidationError{Field: "password", Message: "Password must contain at least one uppercase letter"}
	}
	if !hasLower {
		return ValidationError{Field: "password", Message: "Password must contain at least one lowercase letter"}
	}
	if !hasDigit {
		return ValidationError{Field: "password", Message: "Password must contain at least one digit"}
	}
	if !hasSpecial {
		return ValidationError{Field: "password", Message: "Password must contain at least one special character"}
	}

	return nil
}

// ValidatePhone valida um número de telefone
func ValidatePhone(phone string) error {
	if phone == "" {
		return nil // Phone é opcional
	}

	if !phoneRegex.MatchString(phone) {
		return ValidationError{Field: "phone", Message: "Invalid phone number format"}
	}

	return nil
}

// ValidateUUID valida um UUID
func ValidateUUID(uuid string) error {
	if uuid == "" {
		return ValidationError{Field: "id", Message: "ID is required"}
	}

	if !uuidRegex.MatchString(uuid) {
		return ValidationError{Field: "id", Message: "Invalid ID format"}
	}

	return nil
}

// ValidateName valida um nome
func ValidateName(name string) error {
	if name == "" {
		return ValidationError{Field: "name", Message: "Name is required"}
	}

	if len(name) < 2 {
		return ValidationError{Field: "name", Message: "Name must be at least 2 characters long"}
	}

	if len(name) > 100 {
		return ValidationError{Field: "name", Message: "Name must be at most 100 characters long"}
	}

	return nil
}

// ValidateRole valida um role
func ValidateRole(role string) error {
	validRoles := []string{"user", "admin", "moderator", "super_admin"}

	if role == "" {
		return ValidationError{Field: "role", Message: "Role is required"}
	}

	for _, validRole := range validRoles {
		if role == validRole {
			return nil
		}
	}

	return ValidationError{Field: "role", Message: "Invalid role"}
}

// ValidateStatus valida um status
func ValidateStatus(status string) error {
	validStatuses := []string{"active", "inactive", "pending", "suspended"}

	if status == "" {
		return ValidationError{Field: "status", Message: "Status is required"}
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return ValidationError{Field: "status", Message: "Invalid status"}
}

// SanitizeString limpa e sanitiza uma string
func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

// ValidatePagination valida parâmetros de paginação
func ValidatePagination(page, limit int) error {
	if page < 1 {
		return ValidationError{Field: "page", Message: "Page must be greater than 0"}
	}

	if limit < 1 || limit > 100 {
		return ValidationError{Field: "limit", Message: "Limit must be between 1 and 100"}
	}

	return nil
}

// ValidateEmailList valida uma lista de emails
func ValidateEmailList(emails []string) error {
	if len(emails) == 0 {
		return ValidationError{Field: "emails", Message: "At least one email is required"}
	}

	if len(emails) > 50 {
		return ValidationError{Field: "emails", Message: "Maximum 50 emails allowed"}
	}

	for i, email := range emails {
		if err := ValidateEmail(email); err != nil {
			return ValidationError{Field: "emails", Message: "Invalid email at position " + string(rune(i+1))}
		}
	}

	return nil
}

// ValidateStringLength valida o comprimento de uma string
func ValidateStringLength(field, value string, min, max int) error {
	if value == "" && min > 0 {
		return ValidationError{Field: field, Message: field + " is required"}
	}

	if len(value) < min {
		return ValidationError{Field: field, Message: field + " must be at least " + string(rune(min)) + " characters long"}
	}

	if len(value) > max {
		return ValidationError{Field: field, Message: field + " must be at most " + string(rune(max)) + " characters long"}
	}

	return nil
}

// ValidatePositiveInt valida se um inteiro é positivo
func ValidatePositiveInt(field string, value int) error {
	if value <= 0 {
		return ValidationError{Field: field, Message: field + " must be positive"}
	}
	return nil
}

// ValidateNonNegativeInt valida se um inteiro é não negativo
func ValidateNonNegativeInt(field string, value int) error {
	if value < 0 {
		return ValidationError{Field: field, Message: field + " must be non-negative"}
	}
	return nil
}
