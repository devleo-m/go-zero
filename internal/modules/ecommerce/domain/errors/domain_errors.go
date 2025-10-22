package errors

import "fmt"

// DomainError representa um erro específico do domínio
type DomainError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

// Error implementa a interface error
func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError cria um novo erro de domínio
func NewDomainError(code, message string, details map[string]interface{}) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Erros específicos do domínio
var (
	// Erros de validação
	ErrInvalidEmail    = NewDomainError("INVALID_EMAIL", "email inválido", nil)
	ErrInvalidPassword = NewDomainError("INVALID_PASSWORD", "senha inválida", nil)
	ErrInvalidMoney    = NewDomainError("INVALID_MONEY", "valor monetário inválido", nil)

	// Erros de usuário
	ErrUserNotFound      = NewDomainError("USER_NOT_FOUND", "usuário não encontrado", nil)
	ErrUserAlreadyExists = NewDomainError("USER_ALREADY_EXISTS", "usuário já existe", nil)
	ErrEmailAlreadyInUse = NewDomainError("EMAIL_ALREADY_IN_USE", "email já está em uso", nil)
	ErrUserInactive      = NewDomainError("USER_INACTIVE", "usuário inativo", nil)
	ErrUserNotVerified   = NewDomainError("USER_NOT_VERIFIED", "usuário não verificado", nil)

	// Erros de produto
	ErrProductNotFound     = NewDomainError("PRODUCT_NOT_FOUND", "produto não encontrado", nil)
	ErrProductOutOfStock   = NewDomainError("PRODUCT_OUT_OF_STOCK", "produto fora de estoque", nil)
	ErrInvalidProductPrice = NewDomainError("INVALID_PRODUCT_PRICE", "preço do produto inválido", nil)

	// Erros de pedido
	ErrOrderNotFound    = NewDomainError("ORDER_NOT_FOUND", "pedido não encontrado", nil)
	ErrOrderEmpty       = NewDomainError("ORDER_EMPTY", "pedido vazio", nil)
	ErrOrderAlreadyPaid = NewDomainError("ORDER_ALREADY_PAID", "pedido já foi pago", nil)
	ErrOrderCancelled   = NewDomainError("ORDER_CANCELLED", "pedido cancelado", nil)

	// Erros de validação de negócio
	ErrInsufficientStock  = NewDomainError("INSUFFICIENT_STOCK", "estoque insuficiente", nil)
	ErrInvalidQuantity    = NewDomainError("INVALID_QUANTITY", "quantidade inválida", nil)
	ErrInvalidOrderStatus = NewDomainError("INVALID_ORDER_STATUS", "status do pedido inválido", nil)
)

// IsDomainError verifica se um erro é do domínio
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// GetDomainErrorCode retorna o código do erro de domínio
func GetDomainErrorCode(err error) string {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Code
	}
	return ""
}

// WrapDomainError envolve um erro com contexto adicional
func WrapDomainError(err error, code, message string, details map[string]interface{}) *DomainError {
	return &DomainError{
		Code:    code,
		Message: fmt.Sprintf("%s: %v", message, err),
		Details: details,
	}
}
