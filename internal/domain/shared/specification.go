// internal/domain/shared/specification.go
package shared

import "time"

// Specification representa uma regra de negócio que pode ser reutilizada
type Specification[T any] interface {
	// IsSatisfiedBy verifica se a entidade satisfaz a especificação
	IsSatisfiedBy(entity T) bool

	// ToQueryFilter converte especificação para QueryFilter
	ToQueryFilter() QueryFilter

	// And combina com outra especificação usando AND
	And(other Specification[T]) Specification[T]

	// Or combina com outra especificação usando OR
	Or(other Specification[T]) Specification[T]

	// Not nega a especificação
	Not() Specification[T]
}

// BaseSpecification implementação base
type BaseSpecification[T any] struct {
	filter QueryFilter
}

// NewSpecification cria uma nova especificação
func NewSpecification[T any](filter QueryFilter) Specification[T] {
	return &BaseSpecification[T]{filter: filter}
}

func (s *BaseSpecification[T]) IsSatisfiedBy(entity T) bool {
	// Implementação que valida a entidade contra o filtro
	// (pode ser implementada na infrastructure)
	// Por enquanto retorna true, mas na prática seria implementado
	// com reflection ou interface específica para cada entidade
	return true
}

func (s *BaseSpecification[T]) ToQueryFilter() QueryFilter {
	return s.filter
}

func (s *BaseSpecification[T]) And(other Specification[T]) Specification[T] {
	otherFilter := other.ToQueryFilter()
	newFilter := s.filter
	newFilter.Where = append(newFilter.Where, otherFilter.Where...)
	return NewSpecification[T](newFilter)
}

func (s *BaseSpecification[T]) Or(other Specification[T]) Specification[T] {
	otherFilter := other.ToQueryFilter()
	newFilter := s.filter
	newFilter.Or = append(newFilter.Or, otherFilter.Where)
	return NewSpecification[T](newFilter)
}

func (s *BaseSpecification[T]) Not() Specification[T] {
	// Implementar negação - inverter operadores
	notFilter := s.filter
	for i, condition := range notFilter.Where {
		switch condition.Operator {
		case OpEqual:
			notFilter.Where[i].Operator = OpNotEqual
		case OpNotEqual:
			notFilter.Where[i].Operator = OpEqual
		case OpGreaterThan:
			notFilter.Where[i].Operator = OpLessThanOrEqual
		case OpGreaterThanOrEqual:
			notFilter.Where[i].Operator = OpLessThan
		case OpLessThan:
			notFilter.Where[i].Operator = OpGreaterThanOrEqual
		case OpLessThanOrEqual:
			notFilter.Where[i].Operator = OpGreaterThan
		case OpLike:
			notFilter.Where[i].Operator = OpNotLike
		case OpNotLike:
			notFilter.Where[i].Operator = OpLike
		case OpIn:
			notFilter.Where[i].Operator = OpNotIn
		case OpNotIn:
			notFilter.Where[i].Operator = OpIn
		case OpIsNull:
			notFilter.Where[i].Operator = OpIsNotNull
		case OpIsNotNull:
			notFilter.Where[i].Operator = OpIsNull
		case OpBetween:
			notFilter.Where[i].Operator = OpNotBetween
		case OpNotBetween:
			notFilter.Where[i].Operator = OpBetween
		}
	}
	return NewSpecification[T](notFilter)
}

// ==========================================
// ESPECIFICAÇÕES COMUNS
// ==========================================

// ActiveSpecification especificação para entidades ativas
func ActiveSpecification[T any]() Specification[T] {
	return NewSpecification[T](QueryFilter{
		Where: []Condition{
			{Field: "status", Operator: OpEqual, Value: "active"},
		},
	})
}

// InactiveSpecification especificação para entidades inativas
func InactiveSpecification[T any]() Specification[T] {
	return NewSpecification[T](QueryFilter{
		Where: []Condition{
			{Field: "status", Operator: OpEqual, Value: "inactive"},
		},
	})
}

// CreatedTodaySpecification especificação para entidades criadas hoje
func CreatedTodaySpecification[T any]() Specification[T] {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	return NewSpecification[T](QueryFilter{
		Where: []Condition{
			{Field: "created_at", Operator: OpBetween, Value: []time.Time{today, tomorrow}},
		},
	})
}

// CreatedThisWeekSpecification especificação para entidades criadas esta semana
func CreatedThisWeekSpecification[T any]() Specification[T] {
	now := time.Now()
	weekStart := now.Truncate(24*time.Hour).AddDate(0, 0, -int(now.Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 7)
	return NewSpecification[T](QueryFilter{
		Where: []Condition{
			{Field: "created_at", Operator: OpBetween, Value: []time.Time{weekStart, weekEnd}},
		},
	})
}

// CreatedThisMonthSpecification especificação para entidades criadas este mês
func CreatedThisMonthSpecification[T any]() Specification[T] {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	return NewSpecification[T](QueryFilter{
		Where: []Condition{
			{Field: "created_at", Operator: OpBetween, Value: []time.Time{monthStart, monthEnd}},
		},
	})
}

// RoleSpecification especificação para role específico
func RoleSpecification[T any](role string) Specification[T] {
	return NewSpecification[T](QueryFilter{
		Where: []Condition{
			{Field: "role", Operator: OpEqual, Value: role},
		},
	})
}

// EmailSpecification especificação para email específico
func EmailSpecification[T any](email string) Specification[T] {
	return NewSpecification[T](QueryFilter{
		Where: []Condition{
			{Field: "email", Operator: OpEqual, Value: email},
		},
	})
}

// NameContainsSpecification especificação para nome que contém texto
func NameContainsSpecification[T any](name string) Specification[T] {
	return NewSpecification[T](QueryFilter{
		Where: []Condition{
			{Field: "name", Operator: OpILike, Value: "%" + name + "%"},
		},
	})
}

// ==========================================
// ESPECIFICAÇÕES COMPOSTAS
// ==========================================

// ActiveAdminsSpecification especificação para admins ativos
func ActiveAdminsSpecification[T any]() Specification[T] {
	return ActiveSpecification[T]().And(RoleSpecification[T]("admin"))
}

// ActiveUsersSpecification especificação para usuários ativos
func ActiveUsersSpecification[T any]() Specification[T] {
	return ActiveSpecification[T]().And(RoleSpecification[T]("user"))
}

// PendingActivationSpecification especificação para pendentes de ativação
func PendingActivationSpecification[T any]() Specification[T] {
	return NewSpecification[T](QueryFilter{
		Where: []Condition{
			{Field: "status", Operator: OpEqual, Value: "pending"},
		},
	})
}

// ==========================================
// UTILITÁRIOS
// ==========================================

// CombineSpecifications combina múltiplas especificações com AND
func CombineSpecifications[T any](specs ...Specification[T]) Specification[T] {
	if len(specs) == 0 {
		return NewSpecification[T](QueryFilter{})
	}

	result := specs[0]
	for i := 1; i < len(specs); i++ {
		result = result.And(specs[i])
	}

	return result
}

// AnySpecification combina múltiplas especificações com OR
func AnySpecification[T any](specs ...Specification[T]) Specification[T] {
	if len(specs) == 0 {
		return NewSpecification[T](QueryFilter{})
	}

	result := specs[0]
	for i := 1; i < len(specs); i++ {
		result = result.Or(specs[i])
	}

	return result
}
