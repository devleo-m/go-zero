// internal/domain/shared/aggregation.go
package shared

// AggregationResult representa resultado de agregações
type AggregationResult struct {
	// Count total de registros
	Count int64 `json:"count"`

	// Sum soma de valores
	Sum float64 `json:"sum,omitempty"`

	// Avg média de valores
	Avg float64 `json:"avg,omitempty"`

	// Min valor mínimo
	Min interface{} `json:"min,omitempty"`

	// Max valor máximo
	Max interface{} `json:"max,omitempty"`

	// GroupedResults resultados agrupados
	GroupedResults map[string]*AggregationResult `json:"grouped_results,omitempty"`
}

// NewAggregationResult cria um novo resultado de agregação
func NewAggregationResult() *AggregationResult {
	return &AggregationResult{
		GroupedResults: make(map[string]*AggregationResult),
	}
}

// AddGroupedResult adiciona um resultado agrupado
func (ar *AggregationResult) AddGroupedResult(key string, result *AggregationResult) {
	ar.GroupedResults[key] = result
}

// GetGroupedResult retorna um resultado agrupado
func (ar *AggregationResult) GetGroupedResult(key string) (*AggregationResult, bool) {
	result, exists := ar.GroupedResults[key]
	return result, exists
}

// HasGroupedResults verifica se tem resultados agrupados
func (ar *AggregationResult) HasGroupedResults() bool {
	return len(ar.GroupedResults) > 0
}

// GetGroupKeys retorna as chaves dos grupos
func (ar *AggregationResult) GetGroupKeys() []string {
	keys := make([]string, 0, len(ar.GroupedResults))
	for key := range ar.GroupedResults {
		keys = append(keys, key)
	}
	return keys
}
