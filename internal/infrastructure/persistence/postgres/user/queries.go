package user

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
)

// ==========================================
// QUERIES ESPECÍFICAS SEMPRE PAGINADAS
// ==========================================

// FindActiveUsers busca usuários ativos com paginação
func (r *Repository) FindActiveUsers(ctx context.Context, page, pageSize int) (*shared.PaginatedResult[*user.User], error) {
	filter := shared.NewQueryBuilder().
		WhereEqual("status", "active").
		OrderByDesc("created_at").
		Page(page).
		PageSize(pageSize).
		Build()

	r.logger.Debug("Finding active users", "page", page, "pageSize", pageSize)
	return r.Paginate(ctx, filter)
}

// FindUsersByRole busca usuários por role com paginação
func (r *Repository) FindUsersByRole(ctx context.Context, role string, page, pageSize int) (*shared.PaginatedResult[*user.User], error) {
	filter := shared.NewQueryBuilder().
		WhereEqual("role", role).
		OrderByDesc("created_at").
		Page(page).
		PageSize(pageSize).
		Build()

	r.logger.Debug("Finding users by role", "role", role, "page", page, "pageSize", pageSize)
	return r.Paginate(ctx, filter)
}

// FindUsersByStatus busca usuários por status com paginação
func (r *Repository) FindUsersByStatus(ctx context.Context, status string, page, pageSize int) (*shared.PaginatedResult[*user.User], error) {
	filter := shared.NewQueryBuilder().
		WhereEqual("status", status).
		OrderByDesc("created_at").
		Page(page).
		PageSize(pageSize).
		Build()

	r.logger.Debug("Finding users by status", "status", status, "page", page, "pageSize", pageSize)
	return r.Paginate(ctx, filter)
}

// FindUsersCreatedAfter busca usuários criados após uma data com paginação
func (r *Repository) FindUsersCreatedAfter(ctx context.Context, after time.Time, page, pageSize int) (*shared.PaginatedResult[*user.User], error) {
	filter := shared.NewQueryBuilder().
		WhereGreaterThan("created_at", after).
		OrderByDesc("created_at").
		Page(page).
		PageSize(pageSize).
		Build()

	r.logger.Debug("Finding users created after", "after", after, "page", page, "pageSize", pageSize)
	return r.Paginate(ctx, filter)
}

// FindUsersByEmailDomain busca usuários por domínio de email com paginação
func (r *Repository) FindUsersByEmailDomain(ctx context.Context, domain string, page, pageSize int) (*shared.PaginatedResult[*user.User], error) {
	filter := shared.NewQueryBuilder().
		WhereLike("email", "%@"+domain).
		OrderByDesc("created_at").
		Page(page).
		PageSize(pageSize).
		Build()

	r.logger.Debug("Finding users by email domain", "domain", domain, "page", page, "pageSize", pageSize)
	return r.Paginate(ctx, filter)
}

// SearchUsers busca usuários por termo de busca com paginação
func (r *Repository) SearchUsers(ctx context.Context, searchTerm string, page, pageSize int) (*shared.PaginatedResult[*user.User], error) {
	filter := shared.NewQueryBuilder().
		WhereLike("name", searchTerm).
		OrderByDesc("created_at").
		Page(page).
		PageSize(pageSize).
		Build()

	r.logger.Debug("Searching users", "searchTerm", searchTerm, "page", page, "pageSize", pageSize)
	return r.Paginate(ctx, filter)
}

// ==========================================
// QUERIES DE AUTENTICAÇÃO
// ==========================================

// FindByPasswordResetToken busca usuário por token de reset de senha
func (r *Repository) FindByPasswordResetToken(ctx context.Context, token string) (*user.User, error) {
	var authData UserAuthDataModel

	r.logger.Debug("Finding user by password reset token")

	if err := r.db.WithContext(ctx).
		Where("password_reset_token = ? AND password_reset_expires > ?", token, time.Now()).
		First(&authData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("User not found by password reset token")
			return nil, nil
		}

		r.logger.Error("Failed to find user by password reset token", "error", err)
		return nil, fmt.Errorf("failed to find user by password reset token: %w", err)
	}

	// Buscar o usuário principal
	return r.FindByID(ctx, authData.UserID)
}

// FindByActivationToken busca usuário por token de ativação
func (r *Repository) FindByActivationToken(ctx context.Context, token string) (*user.User, error) {
	var authData UserAuthDataModel

	r.logger.Debug("Finding user by activation token")

	if err := r.db.WithContext(ctx).
		Where("activation_token = ? AND activation_expires > ?", token, time.Now()).
		First(&authData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("User not found by activation token")
			return nil, nil
		}

		r.logger.Error("Failed to find user by activation token", "error", err)
		return nil, fmt.Errorf("failed to find user by activation token: %w", err)
	}

	// Buscar o usuário principal
	return r.FindByID(ctx, authData.UserID)
}

// FindByRefreshToken busca usuário por refresh token
func (r *Repository) FindByRefreshToken(ctx context.Context, token string) (*user.User, error) {
	var authData UserAuthDataModel

	r.logger.Debug("Finding user by refresh token")

	if err := r.db.WithContext(ctx).
		Where("refresh_token = ? AND refresh_expires > ?", token, time.Now()).
		First(&authData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("User not found by refresh token")
			return nil, nil
		}

		r.logger.Error("Failed to find user by refresh token", "error", err)
		return nil, fmt.Errorf("failed to find user by refresh token: %w", err)
	}

	// Buscar o usuário principal
	return r.FindByID(ctx, authData.UserID)
}

// ==========================================
// QUERIES DE ESTATÍSTICAS
// ==========================================

// GetUserStatsByPeriod retorna estatísticas de usuários por período
func (r *Repository) GetUserStatsByPeriod(ctx context.Context, start, end time.Time) (*shared.AggregationResult, error) {
	var stats struct {
		TotalCreated   int64 `json:"total_created"`
		TotalActive    int64 `json:"total_active"`
		TotalInactive  int64 `json:"total_inactive"`
		TotalPending   int64 `json:"total_pending"`
		TotalSuspended int64 `json:"total_suspended"`
		NewUsers       int64 `json:"new_users"`
		ReturningUsers int64 `json:"returning_users"`
	}

	r.logger.Debug("Getting user stats by period", "start", start, "end", end)

	// Total de usuários criados no período
	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&stats.TotalCreated).Error; err != nil {
		r.logger.Error("Failed to count users created in period", "error", err)
		return nil, fmt.Errorf("failed to count users created in period: %w", err)
	}

	// Usuários por status no período
	statusCounts := make(map[string]int64)
	var results []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Select("status, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("status").
		Scan(&results).Error; err != nil {
		r.logger.Error("Failed to count users by status in period", "error", err)
		return nil, fmt.Errorf("failed to count users by status in period: %w", err)
	}

	for _, result := range results {
		statusCounts[result.Status] = result.Count
	}

	stats.TotalActive = statusCounts["active"]
	stats.TotalInactive = statusCounts["inactive"]
	stats.TotalPending = statusCounts["pending"]
	stats.TotalSuspended = statusCounts["suspended"]

	// Usuários que fizeram login no período (aproximação)
	if err := r.db.WithContext(ctx).
		Model(&UserProfileModel{}).
		Where("last_login_at BETWEEN ? AND ?", start, end).
		Count(&stats.ReturningUsers).Error; err != nil {
		r.logger.Error("Failed to count returning users", "error", err)
		return nil, fmt.Errorf("failed to count returning users: %w", err)
	}

	// Novos usuários (criados no período)
	stats.NewUsers = stats.TotalCreated

	r.logger.Debug("User stats by period retrieved",
		"total_created", stats.TotalCreated,
		"new_users", stats.NewUsers,
		"returning_users", stats.ReturningUsers,
	)

	return &shared.AggregationResult{
		Count: stats.TotalCreated,
	}, nil
}

// GetTopUsersByActivity retorna usuários mais ativos
func (r *Repository) GetTopUsersByActivity(ctx context.Context, limit int) ([]*user.User, error) {
	if limit <= 0 {
		limit = 10
	}

	var models []*UserModel

	r.logger.Debug("Getting top users by activity", "limit", limit)

	if err := r.db.WithContext(ctx).
		Table("users u").
		Select("u.*").
		Joins("LEFT JOIN user_profiles p ON u.id = p.user_id").
		Order("p.login_count DESC, u.created_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		r.logger.Error("Failed to get top users by activity", "error", err)
		return nil, fmt.Errorf("failed to get top users by activity: %w", err)
	}

	users, err := ToDomainSlice(models)
	if err != nil {
		r.logger.Error("Failed to convert top users to domain", "error", err)
		return nil, fmt.Errorf("failed to convert top users to domain: %w", err)
	}

	r.logger.Debug("Top users by activity retrieved", "count", len(users))
	return users, nil
}

// ==========================================
// QUERIES DE MANUTENÇÃO
// ==========================================

// FindInactiveUsers busca usuários inativos há mais de X dias
func (r *Repository) FindInactiveUsers(ctx context.Context, days int, page, pageSize int) (*shared.PaginatedResult[*user.User], error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	filter := shared.NewQueryBuilder().
		WhereEqual("status", "inactive").
		WhereLessThan("updated_at", cutoffDate).
		OrderByDesc("updated_at").
		Page(page).
		PageSize(pageSize).
		Build()

	r.logger.Debug("Finding inactive users", "days", days, "cutoff", cutoffDate, "page", page, "pageSize", pageSize)
	return r.Paginate(ctx, filter)
}

// FindUsersWithoutProfile busca usuários sem perfil completo
func (r *Repository) FindUsersWithoutProfile(ctx context.Context, page, pageSize int) (*shared.PaginatedResult[*user.User], error) {
	var models []*UserModel

	r.logger.Debug("Finding users without profile", "page", page, "pageSize", pageSize)

	offset := (page - 1) * pageSize
	if page <= 0 {
		offset = 0
	}

	if err := r.db.WithContext(ctx).
		Table("users u").
		Select("u.*").
		Joins("LEFT JOIN user_profiles p ON u.id = p.user_id").
		Where("p.user_id IS NULL").
		Order("u.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error; err != nil {
		r.logger.Error("Failed to find users without profile", "error", err)
		return nil, fmt.Errorf("failed to find users without profile: %w", err)
	}

	users, err := ToDomainSlice(models)
	if err != nil {
		r.logger.Error("Failed to convert users to domain", "error", err)
		return nil, fmt.Errorf("failed to convert users to domain: %w", err)
	}

	// Contar total para paginação
	var total int64
	if err := r.db.WithContext(ctx).
		Table("users u").
		Joins("LEFT JOIN user_profiles p ON u.id = p.user_id").
		Where("p.user_id IS NULL").
		Count(&total).Error; err != nil {
		r.logger.Error("Failed to count users without profile", "error", err)
		return nil, fmt.Errorf("failed to count users without profile: %w", err)
	}

	// Calcular metadados de paginação
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	return &shared.PaginatedResult[*user.User]{
		Data: users,
		Pagination: shared.PaginationMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			PageSize:    pageSize,
			TotalItems:  total,
			ItemsInPage: len(users),
			HasNext:     hasNext,
			HasPrevious: hasPrev,
		},
	}, nil
}

// ==========================================
// QUERIES DE RELATÓRIOS
// ==========================================

// GetUserGrowthReport retorna relatório de crescimento de usuários
func (r *Repository) GetUserGrowthReport(ctx context.Context, start, end time.Time, groupBy string) (*shared.AggregationResult, error) {
	var results []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	r.logger.Debug("Getting user growth report", "start", start, "end", end, "groupBy", groupBy)

	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Select(fmt.Sprintf("DATE_TRUNC('%s', created_at) as date, COUNT(*) as count", groupBy)).
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Scan(&results).Error; err != nil {
		r.logger.Error("Failed to get user growth report", "error", err)
		return nil, fmt.Errorf("failed to get user growth report: %w", err)
	}

	r.logger.Debug("User growth report retrieved", "data_points", len(results))

	return &shared.AggregationResult{
		Count: int64(len(results)),
	}, nil
}
