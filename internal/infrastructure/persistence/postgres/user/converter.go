package user

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
)

// ==========================================
// CONVERSÃO DOMAIN → MODEL
// ==========================================

// ToModel converte User domain para UserModel
func ToModel(domainUser *user.User) *UserModel {
	if domainUser == nil {
		return nil
	}

	return &UserModel{
		ID:           domainUser.ID,
		CreatedAt:    domainUser.CreatedAt,
		UpdatedAt:    domainUser.UpdatedAt,
		DeletedAt:    convertDeletedAt(domainUser.DeletedAt),
		Name:         domainUser.Name,
		Email:        domainUser.Email.String(),
		PasswordHash: domainUser.Password.String(),
		Phone:        convertPhoneToModel(domainUser.Phone),
		Status:       domainUser.Status.String(),
		Role:         domainUser.Role.String(),
	}
}

// ToProfileModel converte dados de perfil do domain para UserProfileModel
func ToProfileModel(domainUser *user.User) *UserProfileModel {
	if domainUser == nil {
		return nil
	}

	// Aqui você pode adicionar lógica para extrair dados de perfil
	// do domain user se necessário
	return &UserProfileModel{
		UserID: domainUser.ID,
		// Outros campos podem ser preenchidos conforme necessário
	}
}

// ToAuthDataModel converte dados de auth do domain para UserAuthDataModel
func ToAuthDataModel(domainUser *user.User) *UserAuthDataModel {
	if domainUser == nil {
		return nil
	}

	return &UserAuthDataModel{
		UserID: domainUser.ID,
		// Outros campos podem ser preenchidos conforme necessário
	}
}

// ToPreferencesModel converte preferências do domain para UserPreferencesModel
func ToPreferencesModel(domainUser *user.User) *UserPreferencesModel {
	if domainUser == nil {
		return nil
	}

	return &UserPreferencesModel{
		UserID: domainUser.ID,
		// Outros campos podem ser preenchidos conforme necessário
	}
}

// ==========================================
// CONVERSÃO MODEL → DOMAIN
// ==========================================

// ToDomain converte UserModel para User domain
func ToDomain(model *UserModel) (*user.User, error) {
	if model == nil {
		return nil, nil
	}

	// Converter Value Objects
	email, err := convertEmailToDomain(model.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	password := user.NewPasswordFromHash(model.PasswordHash)
	phone := convertPhoneToDomain(model.Phone)
	status := convertStatusToDomain(model.Status)
	role := convertRoleToDomain(model.Role)

	return &user.User{
		BaseEntity: convertBaseEntity(model),
		Name:       model.Name,
		Email:      email,
		Password:   password,
		Phone:      phone,
		Status:     status,
		Role:       role,
	}, nil
}

// ToDomainWithRelations converte UserModel com relacionamentos para User domain
func ToDomainWithRelations(model *UserModel) (*user.User, error) {
	domainUser, err := ToDomain(model)
	if err != nil {
		return nil, err
	}

	if domainUser == nil {
		return nil, nil
	}

	// Aqui você pode adicionar lógica para converter os relacionamentos
	// se necessário no futuro

	return domainUser, nil
}

// ==========================================
// CONVERSÃO EM LOTE
// ==========================================

// ToDomainSlice converte slice de UserModel para slice de User domain
func ToDomainSlice(models []*UserModel) ([]*user.User, error) {
	if len(models) == 0 {
		return []*user.User{}, nil
	}

	users := make([]*user.User, 0, len(models))
	for _, model := range models {
		domainUser, err := ToDomain(model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user %s: %w", model.ID, err)
		}
		users = append(users, domainUser)
	}

	return users, nil
}

// ToModelSlice converte slice de User domain para slice de UserModel
func ToModelSlice(users []*user.User) []*UserModel {
	if len(users) == 0 {
		return []*UserModel{}
	}

	models := make([]*UserModel, 0, len(users))
	for _, user := range users {
		models = append(models, ToModel(user))
	}

	return models
}

// ==========================================
// CONVERSÃO DE QUERY FILTERS
// ==========================================

// QueryFilterToGORM converte shared.QueryFilter para GORM scopes
func QueryFilterToGORM(filter shared.QueryFilter) (func(*gorm.DB) *gorm.DB, error) {
	return func(db *gorm.DB) *gorm.DB {
		query := db

		// Aplicar filtros WHERE
		for _, condition := range filter.Where {
			switch condition.Operator {
			case "=", "eq":
				query = query.Where(fmt.Sprintf("%s = ?", condition.Field), condition.Value)
			case "!=", "ne":
				query = query.Where(fmt.Sprintf("%s != ?", condition.Field), condition.Value)
			case ">", "gt":
				query = query.Where(fmt.Sprintf("%s > ?", condition.Field), condition.Value)
			case ">=", "gte":
				query = query.Where(fmt.Sprintf("%s >= ?", condition.Field), condition.Value)
			case "<", "lt":
				query = query.Where(fmt.Sprintf("%s < ?", condition.Field), condition.Value)
			case "<=", "lte":
				query = query.Where(fmt.Sprintf("%s <= ?", condition.Field), condition.Value)
			case "LIKE", "like":
				query = query.Where(fmt.Sprintf("%s LIKE ?", condition.Field), condition.Value)
			case "ILIKE", "ilike":
				query = query.Where(fmt.Sprintf("%s ILIKE ?", condition.Field), condition.Value)
			case "IN", "in":
				query = query.Where(fmt.Sprintf("%s IN ?", condition.Field), condition.Value)
			case "NOT IN", "not_in":
				query = query.Where(fmt.Sprintf("%s NOT IN ?", condition.Field), condition.Value)
			case "BETWEEN", "between":
				if values, ok := condition.Value.([]interface{}); ok && len(values) == 2 {
					query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", condition.Field), values[0], values[1])
				}
			case "IS NULL", "is_null":
				query = query.Where(fmt.Sprintf("%s IS NULL", condition.Field))
			case "IS NOT NULL", "is_not_null":
				query = query.Where(fmt.Sprintf("%s IS NOT NULL", condition.Field))
			}
		}

		// Aplicar ordenação
		if len(filter.OrderBy) > 0 {
			for _, order := range filter.OrderBy {
				if order.Order == "DESC" || order.Order == "desc" {
					query = query.Order(fmt.Sprintf("%s DESC", order.Field))
				} else {
					query = query.Order(fmt.Sprintf("%s ASC", order.Field))
				}
			}
		}

		// Aplicar paginação
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}

		return query
	}, nil
}

// ==========================================
// MÉTODOS PRIVADOS DE CONVERSÃO
// ==========================================

// convertDeletedAt converte *time.Time para gorm.DeletedAt
func convertDeletedAt(deletedAt *time.Time) gorm.DeletedAt {
	if deletedAt == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Time: *deletedAt, Valid: true}
}

// convertBaseEntity converte campos base do model para shared.BaseEntity
func convertBaseEntity(model *UserModel) shared.BaseEntity {
	base := shared.BaseEntity{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	if model.DeletedAt.Valid {
		base.DeletedAt = &model.DeletedAt.Time
	}

	return base
}

// convertEmailToDomain converte string para user.Email
func convertEmailToDomain(email string) (user.Email, error) {
	return user.NewEmail(email)
}

// convertPhoneToModel converte *user.Phone para *string
func convertPhoneToModel(phone *user.Phone) *string {
	if phone == nil {
		return nil
	}
	str := phone.String()
	return &str
}

// convertPhoneToDomain converte *string para *user.Phone
func convertPhoneToDomain(phone *string) *user.Phone {
	if phone == nil {
		return nil
	}

	p, err := user.NewPhone(*phone)
	if err != nil {
		return nil // ou log error
	}

	return &p
}

// convertStatusToDomain converte string para user.Status
func convertStatusToDomain(status string) user.Status {
	return user.Status(status)
}

// convertRoleToDomain converte string para user.Role
func convertRoleToDomain(role string) user.Role {
	return user.Role(role)
}
