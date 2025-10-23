package user

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
)

// ListUsersUseCase implementa o caso de uso de listagem de usuários
type ListUsersUseCase struct {
	userRepo    user.Repository
	queryHelper *UserQueryHelper
	logger      Logger
}

// NewListUsersUseCase cria uma nova instância do caso de uso
func NewListUsersUseCase(
	userRepo user.Repository,
	queryHelper *UserQueryHelper,
	logger Logger,
) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo:    userRepo,
		queryHelper: queryHelper,
		logger:      logger,
	}
}

// Execute executa o caso de uso de listagem de usuários
func (uc *ListUsersUseCase) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	// 1. VALIDAR E NORMALIZAR INPUT
	if err := uc.validateAndNormalizeInput(&input); err != nil {
		uc.logger.Warn("Invalid input for list users",
			"error", err,
			"input", input,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting list users process",
		"page", input.Page,
		"limit", input.Limit,
		"role", input.Role,
		"status", input.Status,
		"search", input.Search,
	)

	// 2. CONSTRUIR FILTROS DE QUERY
	filter, err := uc.buildQueryFilter(input)
	if err != nil {
		uc.logger.Error("Failed to build query filter",
			"error", err,
			"input", input,
		)
		return nil, fmt.Errorf("failed to build query filter: %w", err)
	}

	uc.logger.Debug("Query filter built successfully",
		"where_clauses", len(filter.Where),
		"order_by_clauses", len(filter.OrderBy),
		"page", filter.Page,
		"limit", filter.Limit,
	)

	// 3. BUSCAR USUÁRIOS USANDO REPOSITORY GENÉRICO
	result, err := uc.userRepo.Paginate(ctx, filter)
	if err != nil {
		uc.logger.Error("Failed to paginate users",
			"error", err,
			"filter", filter,
		)
		return nil, fmt.Errorf("failed to paginate users: %w", err)
	}

	uc.logger.Debug("Users retrieved successfully",
		"total_items", result.Pagination.TotalItems,
		"current_page", result.Pagination.CurrentPage,
		"total_pages", result.Pagination.TotalPages,
		"items_in_page", result.Pagination.ItemsInPage,
	)

	// 4. CONVERTER PARA OUTPUT
	users := make([]UserOutput, len(result.Data))
	for i, domainUser := range result.Data {
		users[i] = ToUserOutput(domainUser)
	}

	// 5. CONSTRUIR PAGINAÇÃO
	pagination := PaginationOutput{
		CurrentPage: result.Pagination.CurrentPage,
		TotalPages:  result.Pagination.TotalPages,
		PageSize:    result.Pagination.PageSize,
		TotalItems:  int(result.Pagination.TotalItems),
		ItemsInPage: result.Pagination.ItemsInPage,
		HasNext:     result.Pagination.CurrentPage < result.Pagination.TotalPages,
		HasPrevious: result.Pagination.CurrentPage > 1,
	}

	// 6. RETORNAR RESULTADO
	output := &ListUsersOutput{
		Users:      users,
		Pagination: pagination,
	}

	uc.logger.Info("List users use case completed successfully",
		"total_users", len(users),
		"page", input.Page,
		"limit", input.Limit,
	)

	return output, nil
}

// validateAndNormalizeInput valida e normaliza os dados de entrada
func (uc *ListUsersUseCase) validateAndNormalizeInput(input *ListUsersInput) error {
	// Normalizar página
	if input.Page <= 0 {
		input.Page = 1
	}

	// Normalizar tamanho da página
	if input.Limit <= 0 {
		input.Limit = 20 // default
	}
	if input.Limit > 100 {
		input.Limit = 100 // máximo
	}

	// Validar role se fornecida
	if input.Role != "" && !user.Role(input.Role).IsValid() {
		return fmt.Errorf("invalid role: %s", input.Role)
	}

	// Validar status se fornecido
	if input.Status != "" && !user.Status(input.Status).IsValid() {
		return fmt.Errorf("invalid status: %s", input.Status)
	}

	return nil
}

// buildQueryFilter constrói os filtros de query
func (uc *ListUsersUseCase) buildQueryFilter(input ListUsersInput) (shared.QueryFilter, error) {
	filter := shared.QueryFilter{
		Page:    input.Page,
		Limit:   input.Limit,
		Where:   []shared.Condition{},
		OrderBy: []shared.OrderBy{},
	}

	// Adicionar filtro de role
	if input.Role != "" {
		filter.Where = append(filter.Where, shared.Condition{
			Field:    "role",
			Operator: shared.OpEqual,
			Value:    input.Role,
		})
	}

	// Adicionar filtro de status
	if input.Status != "" {
		filter.Where = append(filter.Where, shared.Condition{
			Field:    "status",
			Operator: shared.OpEqual,
			Value:    input.Status,
		})
	}

	// Adicionar filtro de busca (busca em nome e email)
	if input.Search != "" {
		searchTerm := "%" + input.Search + "%"
		filter.Where = append(filter.Where, shared.Condition{
			Field:    "name",
			Operator: shared.OpILike,
			Value:    searchTerm,
		})
		filter.Where = append(filter.Where, shared.Condition{
			Field:    "email",
			Operator: shared.OpILike,
			Value:    searchTerm,
		})
	}

	// Adicionar ordenação padrão (por data de criação, mais recentes primeiro)
	filter.OrderBy = append(filter.OrderBy, shared.OrderBy{
		Field: "created_at",
		Order: shared.SortDesc,
	})

	// Adicionar ordenação secundária por nome
	filter.OrderBy = append(filter.OrderBy, shared.OrderBy{
		Field: "name",
		Order: shared.SortAsc,
	})

	uc.logger.Debug("Query filter built",
		"where_clauses", len(filter.Where),
		"order_by_clauses", len(filter.OrderBy),
		"page", filter.Page,
		"limit", filter.Limit,
	)

	return filter, nil
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// GetUsersByRole lista usuários por role
func (uc *ListUsersUseCase) GetUsersByRole(ctx context.Context, role string, page, pageSize int) (*ListUsersOutput, error) {
	input := ListUsersInput{
		Role:  role,
		Page:  page,
		Limit: pageSize,
	}
	return uc.Execute(ctx, input)
}

// GetUsersByStatus lista usuários por status
func (uc *ListUsersUseCase) GetUsersByStatus(ctx context.Context, status string, page, pageSize int) (*ListUsersOutput, error) {
	input := ListUsersInput{
		Status: status,
		Page:   page,
		Limit:  pageSize,
	}
	return uc.Execute(ctx, input)
}

// SearchUsers busca usuários por termo
func (uc *ListUsersUseCase) SearchUsers(ctx context.Context, searchTerm string, page, pageSize int) (*ListUsersOutput, error) {
	input := ListUsersInput{
		Search: searchTerm,
		Page:   page,
		Limit:  pageSize,
	}
	return uc.Execute(ctx, input)
}

// GetActiveUsers lista usuários ativos
func (uc *ListUsersUseCase) GetActiveUsers(ctx context.Context, page, pageSize int) (*ListUsersOutput, error) {
	input := ListUsersInput{
		Status: "active",
		Page:   page,
		Limit:  pageSize,
	}
	return uc.Execute(ctx, input)
}

// GetPendingUsers lista usuários pendentes
func (uc *ListUsersUseCase) GetPendingUsers(ctx context.Context, page, pageSize int) (*ListUsersOutput, error) {
	input := ListUsersInput{
		Status: "pending",
		Page:   page,
		Limit:  pageSize,
	}
	return uc.Execute(ctx, input)
}

// GetUserStats retorna estatísticas dos usuários
func (uc *ListUsersUseCase) GetUserStats(ctx context.Context) (*UserStatsOutput, error) {
	// Usar o helper para obter estatísticas
	stats, err := uc.queryHelper.GetUserStats(ctx)
	if err != nil {
		uc.logger.Error("Failed to get user stats",
			"error", err,
		)
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	uc.logger.Info("User stats retrieved successfully",
		"total", stats.Total,
		"active", stats.Active,
		"pending", stats.Pending,
		"suspended", stats.Suspended,
		"inactive", stats.Inactive,
	)

	return stats, nil
}
