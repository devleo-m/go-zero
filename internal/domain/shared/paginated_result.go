// internal/domain/shared/paginated_result.go
package shared

// PaginatedResult representa resultado paginado genérico
type PaginatedResult[T any] struct {
	// ==========================================
	// DADOS
	// ==========================================

	// Data slice de entidades
	Data []T `json:"data"`

	// ==========================================
	// METADADOS DE PAGINAÇÃO
	// ==========================================

	// Pagination informações de paginação
	Pagination PaginationMeta `json:"pagination"`

	// ==========================================
	// METADADOS DE AGREGAÇÃO (OPCIONAL)
	// ==========================================

	// Aggregations agregações calculadas
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`

	// ==========================================
	// METADADOS DE FILTROS (OPCIONAL)
	// ==========================================

	// AppliedFilters filtros que foram aplicados
	AppliedFilters *QueryFilter `json:"applied_filters,omitempty"`
}

// PaginationMeta metadados de paginação
type PaginationMeta struct {
	// Página atual
	CurrentPage int `json:"current_page"`

	// Total de páginas
	TotalPages int `json:"total_pages"`

	// Itens por página
	PageSize int `json:"page_size"`

	// Total de itens (em todas as páginas)
	TotalItems int64 `json:"total_items"`

	// Número de itens na página atual
	ItemsInPage int `json:"items_in_page"`

	// Tem página anterior?
	HasPrevious bool `json:"has_previous"`

	// Tem próxima página?
	HasNext bool `json:"has_next"`

	// Página anterior (nil se não tiver)
	PreviousPage *int `json:"previous_page,omitempty"`

	// Próxima página (nil se não tiver)
	NextPage *int `json:"next_page,omitempty"`

	// Índice do primeiro item (1-based)
	FirstItemIndex int `json:"first_item_index"`

	// Índice do último item
	LastItemIndex int `json:"last_item_index"`
}

// NewPaginatedResult cria um novo resultado paginado
func NewPaginatedResult[T any](
	data []T,
	total int64,
	page int,
	pageSize int,
) *PaginatedResult[T] {
	totalPages := calculateTotalPages(total, pageSize)
	hasPrevious := page > 1
	hasNext := page < totalPages

	var previousPage, nextPage *int
	if hasPrevious {
		prev := page - 1
		previousPage = &prev
	}
	if hasNext {
		next := page + 1
		nextPage = &next
	}

	firstItemIndex := (page-1)*pageSize + 1
	lastItemIndex := firstItemIndex + len(data) - 1

	return &PaginatedResult[T]{
		Data: data,
		Pagination: PaginationMeta{
			CurrentPage:    page,
			TotalPages:     totalPages,
			PageSize:       pageSize,
			TotalItems:     total,
			ItemsInPage:    len(data),
			HasPrevious:    hasPrevious,
			HasNext:        hasNext,
			PreviousPage:   previousPage,
			NextPage:       nextPage,
			FirstItemIndex: firstItemIndex,
			LastItemIndex:  lastItemIndex,
		},
	}
}

// WithAggregations adiciona agregações ao resultado
func (pr *PaginatedResult[T]) WithAggregations(agg map[string]interface{}) *PaginatedResult[T] {
	pr.Aggregations = agg
	return pr
}

// WithFilters adiciona filtros aplicados ao resultado
func (pr *PaginatedResult[T]) WithFilters(filter *QueryFilter) *PaginatedResult[T] {
	pr.AppliedFilters = filter
	return pr
}

// IsEmpty verifica se o resultado está vazio
func (pr *PaginatedResult[T]) IsEmpty() bool {
	return len(pr.Data) == 0
}

// IsFirstPage verifica se é a primeira página
func (pr *PaginatedResult[T]) IsFirstPage() bool {
	return pr.Pagination.CurrentPage == 1
}

// IsLastPage verifica se é a última página
func (pr *PaginatedResult[T]) IsLastPage() bool {
	return pr.Pagination.CurrentPage == pr.Pagination.TotalPages
}

// GetPageRange retorna o range de páginas para exibir em paginação UI
// Ex: [1, 2, 3, 4, 5] para página 3 de 10
func (pr *PaginatedResult[T]) GetPageRange(maxPages int) []int {
	if maxPages <= 0 {
		maxPages = 5
	}

	total := pr.Pagination.TotalPages
	current := pr.Pagination.CurrentPage

	if total <= maxPages {
		pages := make([]int, total)
		for i := range pages {
			pages[i] = i + 1
		}
		return pages
	}

	half := maxPages / 2
	start := current - half
	end := current + half

	if start < 1 {
		start = 1
		end = maxPages
	}

	if end > total {
		end = total
		start = total - maxPages + 1
	}

	pages := make([]int, end-start+1)
	for i := range pages {
		pages[i] = start + i
	}

	return pages
}

// calculateTotalPages calcula o número total de páginas
func calculateTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}

	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	return pages
}
