// internal/domain/shared/query_filter.go
package shared

// QueryFilter define todos os filtros possíveis para queries
type QueryFilter struct {
	// ==========================================
	// CONDIÇÕES WHERE
	// ==========================================

	// Where são as condições principais
	Where []Condition `json:"where,omitempty"`

	// Or permite condições OR
	Or [][]Condition `json:"or,omitempty"`

	// ==========================================
	// ORDENAÇÃO
	// ==========================================

	// OrderBy define campos para ordenação
	OrderBy []OrderBy `json:"order_by,omitempty"`

	// ==========================================
	// PAGINAÇÃO
	// ==========================================

	// Page número da página (começa em 1)
	Page int `json:"page"`

	// PageSize itens por página
	PageSize int `json:"page_size"`

	// ==========================================
	// RELACIONAMENTOS
	// ==========================================

	// Include relacionamentos a serem carregados
	Include []string `json:"include,omitempty"`

	// ==========================================
	// CAMPOS
	// ==========================================

	// Select campos específicos a serem retornados
	Select []string `json:"select,omitempty"`

	// Omit campos a serem omitidos
	Omit []string `json:"omit,omitempty"`

	// ==========================================
	// SOFT DELETE
	// ==========================================

	// IncludeDeleted inclui entidades soft deleted
	IncludeDeleted bool `json:"include_deleted"`

	// OnlyDeleted retorna apenas entidades soft deleted
	OnlyDeleted bool `json:"only_deleted"`

	// ==========================================
	// AGREGAÇÕES
	// ==========================================

	// GroupBy agrupa resultados
	GroupBy []string `json:"group_by,omitempty"`

	// Having condições para grupos
	Having []Condition `json:"having,omitempty"`

	// ==========================================
	// LIMITAÇÕES
	// ==========================================

	// Limit número máximo de resultados
	Limit int `json:"limit,omitempty"`

	// Offset pular N resultados
	Offset int `json:"offset,omitempty"`
}

// Condition representa uma condição WHERE
type Condition struct {
	// Field nome do campo
	Field string `json:"field"`

	// Operator operador de comparação
	Operator Operator `json:"operator"`

	// Value valor para comparação
	Value interface{} `json:"value"`

	// CaseSensitive se a comparação é case-sensitive
	CaseSensitive bool `json:"case_sensitive"`
}

// Operator define operadores de comparação
type Operator string

const (
	// Operadores de comparação
	OpEqual              Operator = "="
	OpNotEqual           Operator = "!="
	OpGreaterThan        Operator = ">"
	OpGreaterThanOrEqual Operator = ">="
	OpLessThan           Operator = "<"
	OpLessThanOrEqual    Operator = "<="

	// Operadores de string
	OpLike       Operator = "LIKE"
	OpNotLike    Operator = "NOT LIKE"
	OpILike      Operator = "ILIKE" // Case-insensitive LIKE
	OpStartsWith Operator = "STARTS_WITH"
	OpEndsWith   Operator = "ENDS_WITH"
	OpContains   Operator = "CONTAINS"

	// Operadores de array
	OpIn    Operator = "IN"
	OpNotIn Operator = "NOT IN"

	// Operadores de null
	OpIsNull    Operator = "IS NULL"
	OpIsNotNull Operator = "IS NOT NULL"

	// Operadores de range
	OpBetween    Operator = "BETWEEN"
	OpNotBetween Operator = "NOT BETWEEN"
)

// OrderBy define ordenação
type OrderBy struct {
	Field string    `json:"field"`
	Order SortOrder `json:"order"`
}

// SortOrder define direção da ordenação
type SortOrder string

const (
	SortAsc  SortOrder = "ASC"
	SortDesc SortOrder = "DESC"
)

// Validate valida o QueryFilter
func (qf *QueryFilter) Validate() error {
	// Validar página
	if qf.Page < 1 {
		qf.Page = 1
	}

	// Validar tamanho da página
	if qf.PageSize < 1 {
		qf.PageSize = 20
	}

	if qf.PageSize > 100 {
		qf.PageSize = 100
	}

	// Validar condições
	for _, cond := range qf.Where {
		if err := cond.Validate(); err != nil {
			return err
		}
	}

	// Validar ordenação
	for _, order := range qf.OrderBy {
		if err := order.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Condition) Validate() error {
	if c.Field == "" {
		return NewDomainError("INVALID_CONDITION", "field cannot be empty")
	}

	if c.Operator == "" {
		return NewDomainError("INVALID_CONDITION", "operator cannot be empty")
	}

	// Validar operadores que não precisam de valor
	if c.Operator == OpIsNull || c.Operator == OpIsNotNull {
		return nil
	}

	if c.Value == nil {
		return NewDomainError("INVALID_CONDITION", "value cannot be nil")
	}

	return nil
}

func (o *OrderBy) Validate() error {
	if o.Field == "" {
		return NewDomainError("INVALID_ORDER", "field cannot be empty")
	}

	if o.Order != SortAsc && o.Order != SortDesc {
		o.Order = SortAsc
	}

	return nil
}

// GetOffset calcula o offset para paginação
func (qf *QueryFilter) GetOffset() int {
	return (qf.Page - 1) * qf.PageSize
}

// GetLimit retorna o limite de resultados
func (qf *QueryFilter) GetLimit() int {
	if qf.Limit > 0 {
		return qf.Limit
	}
	return qf.PageSize
}
