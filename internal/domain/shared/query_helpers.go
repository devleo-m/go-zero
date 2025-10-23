// internal/domain/shared/query_helpers.go
package shared

import "time"

// QueryBuilder facilita criação de QueryFilter
type QueryBuilder struct {
	filter QueryFilter
}

// NewQueryBuilder cria um novo query builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		filter: QueryFilter{
			Page:     1,
			PageSize: 20,
		},
	}
}

// ==========================================
// CONDIÇÕES WHERE
// ==========================================

// Where adiciona condição WHERE
func (qb *QueryBuilder) Where(field string, operator Operator, value interface{}) *QueryBuilder {
	qb.filter.Where = append(qb.filter.Where, Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return qb
}

// WhereEqual atalho para operador =
func (qb *QueryBuilder) WhereEqual(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, OpEqual, value)
}

// WhereNotEqual atalho para operador !=
func (qb *QueryBuilder) WhereNotEqual(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, OpNotEqual, value)
}

// WhereGreaterThan atalho para operador >
func (qb *QueryBuilder) WhereGreaterThan(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, OpGreaterThan, value)
}

// WhereGreaterThanOrEqual atalho para operador >=
func (qb *QueryBuilder) WhereGreaterThanOrEqual(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, OpGreaterThanOrEqual, value)
}

// WhereLessThan atalho para operador <
func (qb *QueryBuilder) WhereLessThan(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, OpLessThan, value)
}

// WhereLessThanOrEqual atalho para operador <=
func (qb *QueryBuilder) WhereLessThanOrEqual(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, OpLessThanOrEqual, value)
}

// WhereIn atalho para operador IN
func (qb *QueryBuilder) WhereIn(field string, values []interface{}) *QueryBuilder {
	return qb.Where(field, OpIn, values)
}

// WhereNotIn atalho para operador NOT IN
func (qb *QueryBuilder) WhereNotIn(field string, values []interface{}) *QueryBuilder {
	return qb.Where(field, OpNotIn, values)
}

// WhereLike atalho para operador LIKE
func (qb *QueryBuilder) WhereLike(field string, value string) *QueryBuilder {
	return qb.Where(field, OpLike, "%"+value+"%")
}

// WhereILike atalho para operador ILIKE (case-insensitive)
func (qb *QueryBuilder) WhereILike(field string, value string) *QueryBuilder {
	return qb.Where(field, OpILike, "%"+value+"%")
}

// WhereStartsWith atalho para STARTS_WITH
func (qb *QueryBuilder) WhereStartsWith(field string, value string) *QueryBuilder {
	return qb.Where(field, OpStartsWith, value+"%")
}

// WhereEndsWith atalho para ENDS_WITH
func (qb *QueryBuilder) WhereEndsWith(field string, value string) *QueryBuilder {
	return qb.Where(field, OpEndsWith, "%"+value)
}

// WhereContains atalho para CONTAINS
func (qb *QueryBuilder) WhereContains(field string, value string) *QueryBuilder {
	return qb.Where(field, OpContains, "%"+value+"%")
}

// WhereNull atalho para IS NULL
func (qb *QueryBuilder) WhereNull(field string) *QueryBuilder {
	return qb.Where(field, OpIsNull, nil)
}

// WhereNotNull atalho para IS NOT NULL
func (qb *QueryBuilder) WhereNotNull(field string) *QueryBuilder {
	return qb.Where(field, OpIsNotNull, nil)
}

// WhereBetween atalho para BETWEEN
func (qb *QueryBuilder) WhereBetween(field string, start, end interface{}) *QueryBuilder {
	return qb.Where(field, OpBetween, []interface{}{start, end})
}

// WhereNotBetween atalho para NOT BETWEEN
func (qb *QueryBuilder) WhereNotBetween(field string, start, end interface{}) *QueryBuilder {
	return qb.Where(field, OpNotBetween, []interface{}{start, end})
}

// ==========================================
// ORDENAÇÃO
// ==========================================

// OrderBy adiciona ordenação
func (qb *QueryBuilder) OrderBy(field string, order SortOrder) *QueryBuilder {
	qb.filter.OrderBy = append(qb.filter.OrderBy, OrderBy{
		Field: field,
		Order: order,
	})
	return qb
}

// OrderByAsc atalho para ordenação ASC
func (qb *QueryBuilder) OrderByAsc(field string) *QueryBuilder {
	return qb.OrderBy(field, SortAsc)
}

// OrderByDesc atalho para ordenação DESC
func (qb *QueryBuilder) OrderByDesc(field string) *QueryBuilder {
	return qb.OrderBy(field, SortDesc)
}

// ==========================================
// PAGINAÇÃO
// ==========================================

// Page define a página
func (qb *QueryBuilder) Page(page int) *QueryBuilder {
	if page < 1 {
		page = 1
	}
	qb.filter.Page = page
	return qb
}

// PageSize define o tamanho da página
func (qb *QueryBuilder) PageSize(size int) *QueryBuilder {
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	qb.filter.PageSize = size
	return qb
}

// ==========================================
// RELACIONAMENTOS E CAMPOS
// ==========================================

// Include adiciona relacionamentos
func (qb *QueryBuilder) Include(relations ...string) *QueryBuilder {
	qb.filter.Include = append(qb.filter.Include, relations...)
	return qb
}

// Select adiciona campos para selecionar
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.filter.Select = append(qb.filter.Select, fields...)
	return qb
}

// Omit adiciona campos para omitir
func (qb *QueryBuilder) Omit(fields ...string) *QueryBuilder {
	qb.filter.Omit = append(qb.filter.Omit, fields...)
	return qb
}

// ==========================================
// SOFT DELETE
// ==========================================

// IncludeDeleted inclui registros deletados
func (qb *QueryBuilder) IncludeDeleted() *QueryBuilder {
	qb.filter.IncludeDeleted = true
	return qb
}

// OnlyDeleted retorna apenas registros deletados
func (qb *QueryBuilder) OnlyDeleted() *QueryBuilder {
	qb.filter.OnlyDeleted = true
	return qb
}

// ==========================================
// AGREGAÇÕES
// ==========================================

// GroupBy adiciona agrupamento
func (qb *QueryBuilder) GroupBy(fields ...string) *QueryBuilder {
	qb.filter.GroupBy = append(qb.filter.GroupBy, fields...)
	return qb
}

// Having adiciona condição HAVING
func (qb *QueryBuilder) Having(field string, operator Operator, value interface{}) *QueryBuilder {
	qb.filter.Having = append(qb.filter.Having, Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return qb
}

// ==========================================
// LIMITAÇÕES
// ==========================================

// Limit define limite de resultados
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.filter.Limit = limit
	return qb
}

// Offset define offset
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.filter.Offset = offset
	return qb
}

// ==========================================
// MÉTODOS DE CONVENIÊNCIA
// ==========================================

// Active apenas registros ativos (status = 'active')
func (qb *QueryBuilder) Active() *QueryBuilder {
	return qb.WhereEqual("status", "active")
}

// Inactive apenas registros inativos (status = 'inactive')
func (qb *QueryBuilder) Inactive() *QueryBuilder {
	return qb.WhereEqual("status", "inactive")
}

// CreatedToday registros criados hoje
func (qb *QueryBuilder) CreatedToday() *QueryBuilder {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	return qb.WhereBetween("created_at", today, tomorrow)
}

// CreatedThisWeek registros criados esta semana
func (qb *QueryBuilder) CreatedThisWeek() *QueryBuilder {
	now := time.Now()
	weekStart := now.Truncate(24*time.Hour).AddDate(0, 0, -int(now.Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 7)
	return qb.WhereBetween("created_at", weekStart, weekEnd)
}

// CreatedThisMonth registros criados este mês
func (qb *QueryBuilder) CreatedThisMonth() *QueryBuilder {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	return qb.WhereBetween("created_at", monthStart, monthEnd)
}

// ==========================================
// BUILD
// ==========================================

// Build retorna o QueryFilter construído
func (qb *QueryBuilder) Build() QueryFilter {
	// Validar filtro antes de retornar
	qb.filter.Validate()
	return qb.filter
}

// Reset limpa o builder para reutilização
func (qb *QueryBuilder) Reset() *QueryBuilder {
	qb.filter = QueryFilter{
		Page:     1,
		PageSize: 20,
	}
	return qb
}
