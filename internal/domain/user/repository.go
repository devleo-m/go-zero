package user

import (
	"github.com/devleo-m/go-zero/internal/domain/shared"
)

// Repository é a interface do repositório de User
// Estende a interface genérica Repository[T] que fornece todos os métodos CRUD
// e operações avançadas (paginação, agregações, transações, etc)
type Repository interface {
	shared.Repository[*User]

	// ==========================================
	// MÉTODOS ESPECÍFICOS (APENAS SE REALMENTE NECESSÁRIOS)
	// ==========================================
	//
	// Tente sempre usar o genérico primeiro!
	// Adicione métodos específicos APENAS quando:
	// - Query é MUITO complexa e usada constantemente
	// - Envolve múltiplas tabelas (joins complexos)
	// - Tem lógica de negócio dentro da query
	// - Performance crítica (query otimizada manualmente)

	// Exemplo de método específico válido (se necessário):
	// FindUsersWithActiveSubscriptionsAndRecentActivity(
	//     ctx context.Context,
	//     days int,
	// ) ([]*UserWithSubscription, error)
}
