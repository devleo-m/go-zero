// internal/infrastructure/persistence/postgres/user/queries.go
package user

import (
	"context"
	"time"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ==========================================
// QUERIES ESPECÍFICAS OTIMIZADAS
// ==========================================

// FindByEmail busca usuário por email (otimizada com índice)
func (r *Repository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var model UserModel

	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewDomainError("USER_NOT_FOUND", "user not found")
		}
		return nil, err
	}

	return ToDomain(&model)
}

// FindByPhone busca usuário por telefone (otimizada com índice)
func (r *Repository) FindByPhone(ctx context.Context, phone string) (*user.User, error) {
	var model UserModel

	if err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewDomainError("USER_NOT_FOUND", "user not found")
		}
		return nil, err
	}

	return ToDomain(&model)
}

// FindByStatus busca usuários por status (otimizada com índice)
func (r *Repository) FindByStatus(ctx context.Context, status string) ([]*user.User, error) {
	var models []*UserModel

	if err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// FindByRole busca usuários por role (otimizada com índice)
func (r *Repository) FindByRole(ctx context.Context, role string) ([]*user.User, error) {
	var models []*UserModel

	if err := r.db.WithContext(ctx).
		Where("role = ?", role).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// FindActiveUsers busca usuários ativos (otimizada)
func (r *Repository) FindActiveUsers(ctx context.Context) ([]*user.User, error) {
	var models []*UserModel

	if err := r.db.WithContext(ctx).
		Where("status = ?", "active").
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// FindPendingActivation busca usuários pendentes de ativação
func (r *Repository) FindPendingActivation(ctx context.Context) ([]*user.User, error) {
	var models []*UserModel

	if err := r.db.WithContext(ctx).
		Where("status = ? AND activation_token IS NOT NULL", "pending").
		Order("created_at ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// ==========================================
// QUERIES DE ESTATÍSTICAS
// ==========================================

// GetUserStats retorna estatísticas de usuários
func (r *Repository) GetUserStats(ctx context.Context) (*UserStats, error) {
	var stats UserStats

	// Total de usuários
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Usuários ativos
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("status = ?", "active").Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	// Usuários inativos
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("status = ?", "inactive").Count(&stats.InactiveUsers).Error; err != nil {
		return nil, err
	}

	// Usuários pendentes
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("status = ?", "pending").Count(&stats.PendingUsers).Error; err != nil {
		return nil, err
	}

	// Usuários suspensos
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("status = ?", "suspended").Count(&stats.SuspendedUsers).Error; err != nil {
		return nil, err
	}

	// Usuários por role
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Select("role, COUNT(*) as count").
		Group("role").
		Scan(&stats.UsersByRole).Error; err != nil {
		return nil, err
	}

	// Usuários criados hoje
	today := time.Now().Truncate(24 * time.Hour)
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("created_at >= ?", today).
		Count(&stats.UsersCreatedToday).Error; err != nil {
		return nil, err
	}

	// Usuários criados esta semana
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("created_at >= ?", weekStart).
		Count(&stats.UsersCreatedThisWeek).Error; err != nil {
		return nil, err
	}

	// Usuários criados este mês
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("created_at >= ?", monthStart).
		Count(&stats.UsersCreatedThisMonth).Error; err != nil {
		return nil, err
	}

	// Usuários com email verificado
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("email_verified_at IS NOT NULL").
		Count(&stats.VerifiedUsers).Error; err != nil {
		return nil, err
	}

	// Usuários com 2FA habilitado
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("two_factor_enabled = ?", true).
		Count(&stats.TwoFactorUsers).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// UserStats representa estatísticas de usuários
type UserStats struct {
	TotalUsers            int64 `json:"total_users"`
	ActiveUsers           int64 `json:"active_users"`
	InactiveUsers         int64 `json:"inactive_users"`
	PendingUsers          int64 `json:"pending_users"`
	SuspendedUsers        int64 `json:"suspended_users"`
	VerifiedUsers         int64 `json:"verified_users"`
	TwoFactorUsers        int64 `json:"two_factor_users"`
	UsersCreatedToday     int64 `json:"users_created_today"`
	UsersCreatedThisWeek  int64 `json:"users_created_this_week"`
	UsersCreatedThisMonth int64 `json:"users_created_this_month"`
	UsersByRole           []struct {
		Role  string `json:"role"`
		Count int64  `json:"count"`
	} `json:"users_by_role"`
}

// ==========================================
// QUERIES DE BUSCA E FILTROS
// ==========================================

// SearchUsers busca usuários por texto (nome, email)
func (r *Repository) SearchUsers(ctx context.Context, query string, limit int) ([]*user.User, error) {
	var models []*UserModel

	searchQuery := "%" + query + "%"

	if err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery).
		Order("created_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// FindUsersByDateRange busca usuários criados em um período
func (r *Repository) FindUsersByDateRange(ctx context.Context, start, end time.Time) ([]*user.User, error) {
	var models []*UserModel

	if err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// FindUsersByLastLogin busca usuários por último login
func (r *Repository) FindUsersByLastLogin(ctx context.Context, days int) ([]*user.User, error) {
	var models []*UserModel

	cutoff := time.Now().AddDate(0, 0, -days)

	if err := r.db.WithContext(ctx).
		Where("last_login_at < ?", cutoff).
		Order("last_login_at ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// FindUsersWithoutLogin busca usuários que nunca fizeram login
func (r *Repository) FindUsersWithoutLogin(ctx context.Context) ([]*user.User, error) {
	var models []*UserModel

	if err := r.db.WithContext(ctx).
		Where("last_login_at IS NULL").
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// ==========================================
// QUERIES DE TOKENS E RECUPERAÇÃO
// ==========================================

// FindByPasswordResetToken busca usuário por token de reset
func (r *Repository) FindByPasswordResetToken(ctx context.Context, token string) (*user.User, error) {
	var model UserModel

	if err := r.db.WithContext(ctx).
		Where("password_reset_token = ? AND password_reset_expires > ?", token, time.Now()).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewDomainError("INVALID_TOKEN", "invalid or expired token")
		}
		return nil, err
	}

	return ToDomain(&model)
}

// FindByActivationToken busca usuário por token de ativação
func (r *Repository) FindByActivationToken(ctx context.Context, token string) (*user.User, error) {
	var model UserModel

	if err := r.db.WithContext(ctx).
		Where("activation_token = ? AND activation_expires > ?", token, time.Now()).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewDomainError("INVALID_TOKEN", "invalid or expired token")
		}
		return nil, err
	}

	return ToDomain(&model)
}

// ==========================================
// QUERIES DE AUDITORIA
// ==========================================

// FindUsersCreatedBy busca usuários criados por um usuário específico
func (r *Repository) FindUsersCreatedBy(ctx context.Context, createdBy uuid.UUID) ([]*user.User, error) {
	var models []*UserModel

	if err := r.db.WithContext(ctx).
		Where("created_by = ?", createdBy).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// FindUsersUpdatedBy busca usuários atualizados por um usuário específico
func (r *Repository) FindUsersUpdatedBy(ctx context.Context, updatedBy uuid.UUID) ([]*user.User, error) {
	var models []*UserModel

	if err := r.db.WithContext(ctx).
		Where("updated_by = ?", updatedBy).
		Order("updated_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// ==========================================
// QUERIES DE PERFORMANCE
// ==========================================

// FindUsersWithPagination busca usuários com paginação otimizada
func (r *Repository) FindUsersWithPagination(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*shared.PaginatedResult[*user.User], error) {
	var models []*UserModel
	var total int64

	// Construir query base
	query := r.db.WithContext(ctx).Model(&UserModel{})

	// Aplicar filtros
	for field, value := range filters {
		query = query.Where(field+" = ?", value)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Buscar registros
	offset := (page - 1) * pageSize
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error; err != nil {
		return nil, err
	}

	// Converter para domínio
	users, err := ToDomainSlice(models)
	if err != nil {
		return nil, err
	}

	// Criar resultado paginado
	result := shared.NewPaginatedResult(users, total, page, pageSize)

	return result, nil
}

// ==========================================
// QUERIES DE LIMPEZA E MANUTENÇÃO
// ==========================================

// FindExpiredTokens busca tokens expirados para limpeza
func (r *Repository) FindExpiredTokens(ctx context.Context) ([]*user.User, error) {
	var models []*UserModel

	now := time.Now()

	if err := r.db.WithContext(ctx).
		Where("(password_reset_token IS NOT NULL AND password_reset_expires < ?) OR (activation_token IS NOT NULL AND activation_expires < ?)", now, now).
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToDomainSlice(models)
}

// CleanExpiredTokens limpa tokens expirados
func (r *Repository) CleanExpiredTokens(ctx context.Context) (int64, error) {
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("(password_reset_token IS NOT NULL AND password_reset_expires < ?) OR (activation_token IS NOT NULL AND activation_expires < ?)", now, now).
		Updates(map[string]interface{}{
			"password_reset_token":   nil,
			"password_reset_expires": nil,
			"activation_token":       nil,
			"activation_expires":     nil,
		})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
