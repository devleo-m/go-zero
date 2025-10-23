// internal/infrastructure/persistence/postgres/user/converter.go
package user

import (
	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
	"gorm.io/gorm"
)

// ==========================================
// DOMAIN → MODEL
// ==========================================

// ToModel converte entidade de domínio para modelo GORM
func ToModel(domainUser *user.User) *UserModel {
	if domainUser == nil {
		return nil
	}

	model := &UserModel{
		ID:         domainUser.ID,
		CreatedAt:  domainUser.CreatedAt,
		UpdatedAt:  domainUser.UpdatedAt,
		DeletedAt:  gorm.DeletedAt{},
		Name:       domainUser.Name,
		Email:      domainUser.Email.String(),
		Status:     domainUser.Status.String(),
		Role:       domainUser.Role.String(),
		LoginCount: domainUser.LoginCount,
	}

	// Soft delete
	if domainUser.IsDeleted() {
		model.DeletedAt = gorm.DeletedAt{
			Time:  *domainUser.DeletedAt,
			Valid: true,
		}
	}

	// Telefone (opcional)
	if domainUser.Phone != nil {
		phone := domainUser.Phone.String()
		model.Phone = &phone
	}

	// Senha (hash)
	model.PasswordHash = domainUser.Password.String()

	// Último login
	if domainUser.LastLoginAt != nil {
		model.LastLoginAt = domainUser.LastLoginAt
	}

	return model
}

// ==========================================
// MODEL → DOMAIN
// ==========================================

// ToDomain converte modelo GORM para entidade de domínio
func ToDomain(model *UserModel) (*user.User, error) {
	if model == nil {
		return nil, nil
	}

	// Criar email
	email, err := user.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	// Criar senha (assumindo que o hash já está correto)
	password := user.NewPasswordFromHash(model.PasswordHash)

	// Criar status
	status := user.Status(model.Status)

	// Criar role
	role := user.Role(model.Role)

	// Criar entidade base
	baseEntity := shared.BaseEntity{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	// Soft delete
	if model.DeletedAt.Valid {
		baseEntity.DeletedAt = &model.DeletedAt.Time
	}

	// Criar usuário
	domainUser := &user.User{
		BaseEntity: baseEntity,
		Name:       model.Name,
		Email:      email,
		Password:   password,
		Status:     status,
		Role:       role,
		LoginCount: model.LoginCount,
	}

	// Telefone (opcional)
	if model.Phone != nil {
		phone, err := user.NewPhone(*model.Phone)
		if err != nil {
			return nil, err
		}
		domainUser.Phone = &phone
	}

	// Último login
	if model.LastLoginAt != nil {
		domainUser.LastLoginAt = model.LastLoginAt
	}

	return domainUser, nil
}

// ==========================================
// CONVERSÃO EM LOTE
// ==========================================

// ToDomainSlice converte slice de modelos para slice de entidades
func ToDomainSlice(models []*UserModel) ([]*user.User, error) {
	if len(models) == 0 {
		return []*user.User{}, nil
	}

	users := make([]*user.User, 0, len(models))
	for _, model := range models {
		user, err := ToDomain(model)
		if err != nil {
			return nil, err
		}
		if user != nil {
			users = append(users, user)
		}
	}

	return users, nil
}

// ToModelSlice converte slice de entidades para slice de modelos
func ToModelSlice(users []*user.User) []*UserModel {
	if len(users) == 0 {
		return []*UserModel{}
	}

	models := make([]*UserModel, 0, len(users))
	for _, user := range users {
		model := ToModel(user)
		if model != nil {
			models = append(models, model)
		}
	}

	return models
}

// ==========================================
// CONVERSÃO DE FILTROS
// ==========================================

// QueryFilterToGORM converte QueryFilter para condições GORM
func QueryFilterToGORM(filter shared.QueryFilter) (func(*gorm.DB) *gorm.DB, error) {
	return func(db *gorm.DB) *gorm.DB {
		// Aplicar condições WHERE
		for _, condition := range filter.Where {
			db = applyCondition(db, condition)
		}

		// Aplicar condições OR
		if len(filter.Or) > 0 {
			for _, orGroup := range filter.Or {
				db = db.Or(func(subQuery *gorm.DB) *gorm.DB {
					for _, condition := range orGroup {
						subQuery = applyCondition(subQuery, condition)
					}
					return subQuery
				})
			}
		}

		// Aplicar ordenação
		for _, order := range filter.OrderBy {
			db = db.Order(order.Field + " " + string(order.Order))
		}

		// Aplicar paginação
		if filter.Page > 0 && filter.PageSize > 0 {
			offset := (filter.Page - 1) * filter.PageSize
			db = db.Offset(offset).Limit(filter.PageSize)
		}

		// Aplicar limite
		if filter.Limit > 0 {
			db = db.Limit(filter.Limit)
		}

		// Aplicar offset
		if filter.Offset > 0 {
			db = db.Offset(filter.Offset)
		}

		// Aplicar relacionamentos
		for _, relation := range filter.Include {
			db = db.Preload(relation)
		}

		// Aplicar seleção de campos
		if len(filter.Select) > 0 {
			db = db.Select(filter.Select)
		}

		// Aplicar omissão de campos
		if len(filter.Omit) > 0 {
			db = db.Omit(filter.Omit...)
		}

		// Aplicar soft delete
		if filter.OnlyDeleted {
			db = db.Unscoped().Where("deleted_at IS NOT NULL")
		} else if !filter.IncludeDeleted {
			db = db.Where("deleted_at IS NULL")
		}

		return db
	}, nil
}

// applyCondition aplica uma condição individual
func applyCondition(db *gorm.DB, condition shared.Condition) *gorm.DB {
	switch condition.Operator {
	case shared.OpEqual:
		return db.Where(condition.Field+" = ?", condition.Value)
	case shared.OpNotEqual:
		return db.Where(condition.Field+" != ?", condition.Value)
	case shared.OpGreaterThan:
		return db.Where(condition.Field+" > ?", condition.Value)
	case shared.OpGreaterThanOrEqual:
		return db.Where(condition.Field+" >= ?", condition.Value)
	case shared.OpLessThan:
		return db.Where(condition.Field+" < ?", condition.Value)
	case shared.OpLessThanOrEqual:
		return db.Where(condition.Field+" <= ?", condition.Value)
	case shared.OpLike:
		return db.Where(condition.Field+" LIKE ?", condition.Value)
	case shared.OpNotLike:
		return db.Where(condition.Field+" NOT LIKE ?", condition.Value)
	case shared.OpILike:
		return db.Where(condition.Field+" ILIKE ?", condition.Value)
	case shared.OpIn:
		return db.Where(condition.Field+" IN ?", condition.Value)
	case shared.OpNotIn:
		return db.Where(condition.Field+" NOT IN ?", condition.Value)
	case shared.OpIsNull:
		return db.Where(condition.Field + " IS NULL")
	case shared.OpIsNotNull:
		return db.Where(condition.Field + " IS NOT NULL")
	case shared.OpBetween:
		if values, ok := condition.Value.([]interface{}); ok && len(values) == 2 {
			return db.Where(condition.Field+" BETWEEN ? AND ?", values[0], values[1])
		}
		return db
	case shared.OpNotBetween:
		if values, ok := condition.Value.([]interface{}); ok && len(values) == 2 {
			return db.Where(condition.Field+" NOT BETWEEN ? AND ?", values[0], values[1])
		}
		return db
	default:
		return db
	}
}
