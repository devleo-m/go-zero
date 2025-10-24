package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Params representa os parâmetros de paginação
type Params struct {
	Page  int
	Limit int
	Sort  string
	Order string
}

// Result representa o resultado paginado
type Result struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// ParseFromQuery extrai parâmetros de paginação da query string
func ParseFromQuery(c *gin.Context) *Params {
	page := parseInt(c.Query("page"), 1)
	limit := parseInt(c.Query("limit"), 10)
	sort := c.Query("sort")
	order := c.Query("order")

	// Validações básicas
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	return &Params{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}
}

// Offset calcula o offset baseado na página e limite
func (p *Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

// NewResult cria um novo resultado paginado
func NewResult(data interface{}, total int64, params *Params) *Result {
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}

	return &Result{
		Data:       data,
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}
}

// parseInt converte string para int com valor padrão
func parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return val
}

// ValidateParams valida os parâmetros de paginação
func ValidateParams(params *Params) error {
	if params.Page < 1 {
		return &ValidationError{Field: "page", Message: "Page must be greater than 0"}
	}

	if params.Limit < 1 {
		return &ValidationError{Field: "limit", Message: "Limit must be greater than 0"}
	}

	if params.Limit > 100 {
		return &ValidationError{Field: "limit", Message: "Limit must be at most 100"}
	}

	if params.Order != "" && params.Order != "asc" && params.Order != "desc" {
		return &ValidationError{Field: "order", Message: "Order must be 'asc' or 'desc'"}
	}

	return nil
}

// ValidationError representa um erro de validação
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// GetSortClause retorna a cláusula ORDER BY para SQL
func (p *Params) GetSortClause() string {
	if p.Sort == "" {
		return ""
	}

	order := "ASC"
	if p.Order == "desc" {
		order = "DESC"
	}

	return p.Sort + " " + order
}

// GetGORMOrder retorna a string de ordenação para GORM
func (p *Params) GetGORMOrder() string {
	if p.Sort == "" {
		return ""
	}

	order := "asc"
	if p.Order == "desc" {
		order = "desc"
	}

	return p.Sort + " " + order
}
