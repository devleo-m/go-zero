package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wrapper para o validator com validações customizadas
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator cria uma nova instância do CustomValidator
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Registrar validações customizadas
	registerCustomValidations(v)

	// Configurar tags customizadas
	configureCustomTags(v)

	return &CustomValidator{
		validator: v,
	}
}

// ValidateStruct valida uma struct
func (cv *CustomValidator) ValidateStruct(s interface{}) error {
	return cv.validator.Struct(s)
}

// ValidateVar valida uma variável
func (cv *CustomValidator) ValidateVar(field interface{}, tag string) error {
	return cv.validator.Var(field, tag)
}

// registerCustomValidations registra validações customizadas
func registerCustomValidations(v *validator.Validate) {
	// Validação de CPF
	v.RegisterValidation("cpf", validateCPF)

	// Validação de CNPJ
	v.RegisterValidation("cnpj", validateCNPJ)

	// Validação de telefone brasileiro
	v.RegisterValidation("phone_br", validatePhoneBR)

	// Validação de CEP brasileiro
	v.RegisterValidation("cep", validateCEP)

	// Validação de senha forte
	v.RegisterValidation("strong_password", validateStrongPassword)

	// Validação de role válido
	v.RegisterValidation("valid_role", validateRole)

	// Validação de status válido
	v.RegisterValidation("valid_status", validateStatus)

	// Validação de UUID v4
	v.RegisterValidation("uuid4", validateUUID4)
}

// configureCustomTags configura tags customizadas
func configureCustomTags(v *validator.Validate) {
	// Configurar tag de nome do campo
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ==========================================
// CUSTOM VALIDATIONS
// ==========================================

// validateCPF valida CPF brasileiro
func validateCPF(fl validator.FieldLevel) bool {
	cpf := fl.Field().String()

	// Remove caracteres não numéricos
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")
	cpf = strings.ReplaceAll(cpf, " ", "")

	// Verificar se tem 11 dígitos
	if len(cpf) != 11 {
		return false
	}

	// Verificar se todos os dígitos são iguais
	if strings.Count(cpf, string(cpf[0])) == 11 {
		return false
	}

	// Validar primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	firstDigit := 0
	if remainder >= 2 {
		firstDigit = 11 - remainder
	}

	if int(cpf[9]-'0') != firstDigit {
		return false
	}

	// Validar segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	secondDigit := 0
	if remainder >= 2 {
		secondDigit = 11 - remainder
	}

	return int(cpf[10]-'0') == secondDigit
}

// validateCNPJ valida CNPJ brasileiro
func validateCNPJ(fl validator.FieldLevel) bool {
	cnpj := fl.Field().String()

	// Remove caracteres não numéricos
	cnpj = strings.ReplaceAll(cnpj, ".", "")
	cnpj = strings.ReplaceAll(cnpj, "/", "")
	cnpj = strings.ReplaceAll(cnpj, "-", "")
	cnpj = strings.ReplaceAll(cnpj, " ", "")

	// Verificar se tem 14 dígitos
	if len(cnpj) != 14 {
		return false
	}

	// Verificar se todos os dígitos são iguais
	if strings.Count(cnpj, string(cnpj[0])) == 14 {
		return false
	}

	// Validar primeiro dígito verificador
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(cnpj[i] - '0')
		sum += digit * weights1[i]
	}
	remainder := sum % 11
	firstDigit := 0
	if remainder >= 2 {
		firstDigit = 11 - remainder
	}

	if int(cnpj[12]-'0') != firstDigit {
		return false
	}

	// Validar segundo dígito verificador
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 13; i++ {
		digit := int(cnpj[i] - '0')
		sum += digit * weights2[i]
	}
	remainder = sum % 11
	secondDigit := 0
	if remainder >= 2 {
		secondDigit = 11 - remainder
	}

	return int(cnpj[13]-'0') == secondDigit
}

// validatePhoneBR valida telefone brasileiro
func validatePhoneBR(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Remove caracteres não numéricos
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Verificar se tem 10 ou 11 dígitos
	if len(phone) != 10 && len(phone) != 11 {
		return false
	}

	// Verificar se começa com DDD válido (11-99)
	ddd := phone[:2]
	if ddd < "11" || ddd > "99" {
		return false
	}

	// Verificar se o terceiro dígito é 9 (celular) ou diferente de 9 (fixo)
	if len(phone) == 11 {
		if phone[2] != '9' {
			return false
		}
	}

	return true
}

// validateCEP valida CEP brasileiro
func validateCEP(fl validator.FieldLevel) bool {
	cep := fl.Field().String()

	// Remove caracteres não numéricos
	cep = strings.ReplaceAll(cep, "-", "")
	cep = strings.ReplaceAll(cep, " ", "")

	// Verificar se tem 8 dígitos
	if len(cep) != 8 {
		return false
	}

	// Verificar se todos são dígitos
	for _, char := range cep {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// validateStrongPassword valida senha forte
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Mínimo 8 caracteres
	if len(password) < 8 {
		return false
	}

	// Máximo 128 caracteres
	if len(password) > 128 {
		return false
	}

	// Deve conter pelo menos uma letra minúscula
	hasLower := false
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return false
	}

	// Deve conter pelo menos uma letra maiúscula
	hasUpper := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return false
	}

	// Deve conter pelo menos um número
	hasNumber := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		return false
	}

	// Deve conter pelo menos um caractere especial
	hasSpecial := false
	specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		return false
	}

	return true
}

// validateRole valida role de usuário
func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	validRoles := []string{"admin", "manager", "user", "guest"}

	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}

	return false
}

// validateStatus valida status de usuário
func validateStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"active", "inactive", "pending", "suspended"}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}

	return false
}

// validateUUID4 valida UUID v4
func validateUUID4(fl validator.FieldLevel) bool {
	uuid := fl.Field().String()

	// Verificar formato básico
	if len(uuid) != 36 {
		return false
	}

	// Verificar se tem hífens nas posições corretas
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return false
	}

	// Verificar se é UUID v4 (terceiro grupo começa com 4)
	if uuid[14] != '4' {
		return false
	}

	// Verificar se o quarto grupo começa com 8, 9, A ou B
	fourthGroup := uuid[19]
	if fourthGroup != '8' && fourthGroup != '9' && fourthGroup != 'a' && fourthGroup != 'b' &&
		fourthGroup != 'A' && fourthGroup != 'B' {
		return false
	}

	return true
}

// ==========================================
// HELPER FUNCTIONS
// ==========================================

// FormatValidationErrors formata erros de validação para exibição
func FormatValidationErrors(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, getValidationMessage(e))
		}
	} else {
		errors = append(errors, err.Error())
	}

	return errors
}

// getValidationMessage obtém mensagem de validação personalizada
func getValidationMessage(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "cpf":
		return fmt.Sprintf("%s must be a valid CPF", field)
	case "cnpj":
		return fmt.Sprintf("%s must be a valid CNPJ", field)
	case "phone_br":
		return fmt.Sprintf("%s must be a valid Brazilian phone number", field)
	case "cep":
		return fmt.Sprintf("%s must be a valid Brazilian postal code", field)
	case "strong_password":
		return fmt.Sprintf("%s must be a strong password (8+ chars, upper, lower, number, special)", field)
	case "valid_role":
		return fmt.Sprintf("%s must be a valid role (admin, manager, user, guest)", field)
	case "valid_status":
		return fmt.Sprintf("%s must be a valid status (active, inactive, pending, suspended)", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// ==========================================
// VALIDATION HELPERS
// ==========================================

// ValidateStructWithCustomMessages valida struct com mensagens customizadas
func ValidateStructWithCustomMessages(v *CustomValidator, s interface{}) map[string]string {
	errors := make(map[string]string)

	if err := v.ValidateStruct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				field := e.Field()
				message := getValidationMessage(e)
				errors[field] = message
			}
		}
	}

	return errors
}

// ValidateRequiredFields valida campos obrigatórios
func ValidateRequiredFields(fields map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	for field, value := range fields {
		if value == nil || value == "" {
			errors[field] = fmt.Sprintf("%s is required", field)
		}
	}

	return errors
}
