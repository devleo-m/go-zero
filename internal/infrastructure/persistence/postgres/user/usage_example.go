// internal/infrastructure/persistence/postgres/user/usage_example.go
package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/devleo-m/go-zero/internal/domain/shared"
)

// ExemploUsoRepository demonstra como usar o User Repository
func ExemploUsoRepository() {
	// Simulando contexto e banco
	// ctx := context.Background()
	// db := setupDatabase() // Na prática seria injetado
	// userRepo := NewRepository(db)

	fmt.Println("🚀 EXEMPLOS DE USO DO USER REPOSITORY")
	fmt.Println(strings.Repeat("=", 50))

	// ==========================================
	// 1. CRUD BÁSICO
	// ==========================================
	fmt.Println("\n📝 1. CRUD BÁSICO")

	// Criar usuário
	fmt.Println("   Criando usuário...")
	// email, _ := user.NewEmail("joao@example.com")
	// password, _ := user.NewPassword("senha123")
	// role := user.RoleUser
	// newUser, _ := user.NewUser("João Silva", email, password, role)
	// err := userRepo.Create(ctx, newUser)

	// Buscar por ID
	fmt.Println("   Buscando usuário por ID...")
	// userID := uuid.New()
	// foundUser, err := userRepo.FindByID(ctx, userID)

	// Atualizar usuário
	fmt.Println("   Atualizando usuário...")
	// foundUser.Name = "João Silva Santos"
	// err = userRepo.Update(ctx, foundUser.ID, foundUser)

	// Deletar usuário (soft delete)
	fmt.Println("   Deletando usuário...")
	// err = userRepo.Delete(ctx, foundUser.ID)

	// ==========================================
	// 2. QUERY BUILDER - FACILITA MUITO!
	// ==========================================
	fmt.Println("\n🔍 2. QUERY BUILDER")

	// Buscar usuários ativos
	fmt.Println("   Buscando usuários ativos...")
	activeFilter := shared.NewQueryBuilder().
		WhereEqual("status", "active").
		OrderByDesc("created_at").
		Page(1).
		PageSize(20).
		Build()

	// Buscar usuários por role
	fmt.Println("   Buscando usuários por role...")
	adminFilter := shared.NewQueryBuilder().
		WhereEqual("role", "admin").
		WhereEqual("status", "active").
		OrderByAsc("name").
		Build()

	// Buscar usuários criados hoje
	fmt.Println("   Buscando usuários criados hoje...")
	todayFilter := shared.NewQueryBuilder().
		CreatedToday().
		OrderByDesc("created_at").
		Build()

	// Buscar usuários por nome (LIKE)
	fmt.Println("   Buscando usuários por nome...")
	searchFilter := shared.NewQueryBuilder().
		WhereILike("name", "joão").
		WhereEqual("status", "active").
		Build()

	// Buscar usuários com múltiplas condições
	fmt.Println("   Buscando usuários com múltiplas condições...")
	complexFilter := shared.NewQueryBuilder().
		WhereEqual("status", "active").
		WhereIn("role", []interface{}{"admin", "manager", "user"}).
		WhereBetween("created_at", time.Now().AddDate(0, -1, 0), time.Now()).
		OrderByDesc("created_at").
		Page(1).
		PageSize(50).
		Build()

	_ = activeFilter
	_ = adminFilter
	_ = todayFilter
	_ = searchFilter
	_ = complexFilter

	// ==========================================
	// 3. SPECIFICATION PATTERN - REUTILIZAÇÃO!
	// ==========================================
	fmt.Println("\n🎯 3. SPECIFICATION PATTERN")

	// Usar especificações reutilizáveis
	fmt.Println("   Usando especificações...")
	// activeUsers := shared.ActiveSpecification[user.User]()
	// adminUsers := shared.RoleSpecification[user.User]("admin")
	// activeAdmins := activeUsers.And(adminUsers)

	// Buscar com especificações
	fmt.Println("   Buscando com especificações...")
	// users, err := userRepo.FindMany(ctx, activeAdmins.ToQueryFilter())

	// Combinar especificações
	fmt.Println("   Combinando especificações...")
	// thisWeek := shared.CreatedThisWeekSpecification[user.User]()
	// activeThisWeek := activeUsers.And(thisWeek)
	// users, err = userRepo.FindMany(ctx, activeThisWeek.ToQueryFilter())

	// ==========================================
	// 4. PAGINAÇÃO PROFISSIONAL
	// ==========================================
	fmt.Println("\n📄 4. PAGINAÇÃO PROFISSIONAL")

	// Buscar com paginação
	fmt.Println("   Buscando com paginação...")
	paginatedFilter := shared.NewQueryBuilder().
		WhereEqual("status", "active").
		OrderByDesc("created_at").
		Page(2).
		PageSize(20).
		Build()

	// result, err := userRepo.Paginate(ctx, paginatedFilter)
	// if err == nil {
	//     fmt.Printf("   Página %d de %d (Total: %d usuários)\n",
	//         result.Pagination.CurrentPage,
	//         result.Pagination.TotalPages,
	//         result.Pagination.TotalItems)
	//     fmt.Printf("   Itens na página: %d\n", result.Pagination.ItemsInPage)
	//     fmt.Printf("   Tem página anterior: %v\n", result.Pagination.HasPrevious)
	//     fmt.Printf("   Tem próxima página: %v\n", result.Pagination.HasNext)
	// }

	_ = paginatedFilter

	// ==========================================
	// 5. AGREGAÇÕES E ESTATÍSTICAS
	// ==========================================
	fmt.Println("\n📊 5. AGREGAÇÕES E ESTATÍSTICAS")

	// Contar usuários
	fmt.Println("   Contando usuários...")
	// count, err := userRepo.Count(ctx, shared.QueryFilter{
	//     Where: []shared.Condition{
	//         {Field: "status", Operator: shared.OpEqual, Value: "active"},
	//     },
	// })

	// Verificar se existe
	fmt.Println("   Verificando se existe...")
	// exists, err := userRepo.Exists(ctx, shared.QueryFilter{
	//     Where: []shared.Condition{
	//         {Field: "email", Operator: shared.OpEqual, Value: "joao@example.com"},
	//     },
	// })

	// Estatísticas completas
	fmt.Println("   Obtendo estatísticas...")
	// stats, err := userRepo.GetUserStats(ctx)
	// if err == nil {
	//     fmt.Printf("   Total de usuários: %d\n", stats.TotalUsers)
	//     fmt.Printf("   Usuários ativos: %d\n", stats.ActiveUsers)
	//     fmt.Printf("   Usuários inativos: %d\n", stats.InactiveUsers)
	//     fmt.Printf("   Usuários pendentes: %d\n", stats.PendingUsers)
	//     fmt.Printf("   Usuários criados hoje: %d\n", stats.UsersCreatedToday)
	// }

	// ==========================================
	// 6. OPERAÇÕES EM LOTE
	// ==========================================
	fmt.Println("\n🔄 6. OPERAÇÕES EM LOTE")

	// Update em lote
	fmt.Println("   Atualizando usuários em lote...")
	// affected, err := userRepo.UpdateMany(ctx,
	//     shared.QueryFilter{
	//         Where: []shared.Condition{
	//             {Field: "status", Operator: shared.OpEqual, Value: "pending"},
	//         },
	//     },
	//     map[string]interface{}{
	//         "status": "active",
	//     },
	// )
	// fmt.Printf("   %d usuários foram ativados\n", affected)

	// Delete em lote
	fmt.Println("   Deletando usuários em lote...")
	// cutoffDate := time.Now().AddDate(0, -6, 0) // 6 meses atrás
	// affected, err = userRepo.DeleteMany(ctx, shared.QueryFilter{
	//     Where: []shared.Condition{
	//         {Field: "created_at", Operator: shared.OpLessThan, Value: cutoffDate},
	//         {Field: "status", Operator: shared.OpEqual, Value: "inactive"},
	//     },
	// })
	// fmt.Printf("   %d usuários foram removidos\n", affected)

	// ==========================================
	// 7. TRANSAÇÕES
	// ==========================================
	fmt.Println("\n💼 7. TRANSAÇÕES")

	// Executar em transação
	fmt.Println("   Executando transação...")
	// err = userRepo.WithTransaction(ctx, func(ctx context.Context) error {
	//     // Criar usuário
	//     // if err := userRepo.Create(ctx, user1); err != nil {
	//     //     return err // Rollback automático
	//     // }
	//
	//     // Atualizar usuário
	//     // if err := userRepo.Update(ctx, user2.ID, user2); err != nil {
	//     //     return err // Rollback automático
	//     // }
	//
	//     return nil // Commit automático
	// })

	// ==========================================
	// 8. QUERIES ESPECÍFICAS OTIMIZADAS
	// ==========================================
	fmt.Println("\n⚡ 8. QUERIES ESPECÍFICAS")

	// Buscar por email (otimizada)
	fmt.Println("   Buscando por email...")
	// user, err := userRepo.FindByEmail(ctx, "joao@example.com")

	// Buscar por telefone (otimizada)
	fmt.Println("   Buscando por telefone...")
	// user, err = userRepo.FindByPhone(ctx, "+5511999999999")

	// Buscar usuários ativos (otimizada)
	fmt.Println("   Buscando usuários ativos...")
	// users, err := userRepo.FindActiveUsers(ctx)

	// Buscar usuários pendentes
	fmt.Println("   Buscando usuários pendentes...")
	// users, err = userRepo.FindPendingActivation(ctx)

	// Buscar por token de reset
	fmt.Println("   Buscando por token de reset...")
	// user, err = userRepo.FindByPasswordResetToken(ctx, "token123")

	// Buscar por token de ativação
	fmt.Println("   Buscando por token de ativação...")
	// user, err = userRepo.FindByActivationToken(ctx, "token456")

	// Buscar usuários sem login
	fmt.Println("   Buscando usuários sem login...")
	// users, err = userRepo.FindUsersWithoutLogin(ctx)

	// Buscar usuários por período
	fmt.Println("   Buscando usuários por período...")
	// start := time.Now().AddDate(0, -1, 0) // 1 mês atrás
	// end := time.Now()
	// users, err = userRepo.FindUsersByDateRange(ctx, start, end)

	// Buscar usuários por último login
	fmt.Println("   Buscando usuários por último login...")
	// users, err = userRepo.FindUsersByLastLogin(ctx, 30) // 30 dias

	// Buscar usuários por texto
	fmt.Println("   Buscando usuários por texto...")
	// users, err = userRepo.SearchUsers(ctx, "joão", 10)

	// ==========================================
	// 9. QUERIES AVANÇADAS
	// ==========================================
	fmt.Println("\n🔬 9. QUERIES AVANÇADAS")

	// Distinct - roles únicos
	fmt.Println("   Buscando roles únicos...")
	// roles, err := userRepo.Distinct(ctx, "role", shared.QueryFilter{})

	// GroupBy - agrupar por role
	fmt.Println("   Agrupando por role...")
	// groups, err := userRepo.GroupBy(ctx, "role", shared.QueryFilter{})

	// FindFirst - primeiro usuário
	fmt.Println("   Buscando primeiro usuário...")
	// firstUser, err := userRepo.FindFirst(ctx, shared.QueryFilter{
	//     OrderBy: []shared.OrderBy{
	//         {Field: "created_at", Order: shared.SortAsc},
	//     },
	// })

	// FindLast - último usuário
	fmt.Println("   Buscando último usuário...")
	// lastUser, err := userRepo.FindLast(ctx, shared.QueryFilter{
	//     OrderBy: []shared.OrderBy{
	//         {Field: "created_at", Order: shared.SortDesc},
	//     },
	// })

	// ==========================================
	// 10. MANUTENÇÃO E LIMPEZA
	// ==========================================
	fmt.Println("\n🧹 10. MANUTENÇÃO E LIMPEZA")

	// Limpar tokens expirados
	fmt.Println("   Limpando tokens expirados...")
	// affected, err := userRepo.CleanExpiredTokens(ctx)
	// fmt.Printf("   %d tokens foram limpos\n", affected)

	// Buscar tokens expirados
	fmt.Println("   Buscando tokens expirados...")
	// users, err := userRepo.FindExpiredTokens(ctx)

	fmt.Println("\n✅ Exemplos concluídos com sucesso!")
}

// ExemploComparacaoAntesDepois demonstra a diferença entre os padrões
func ExemploComparacaoAntesDepois() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 COMPARAÇÃO: ANTES vs AGORA")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n❌ ANTES (Repository específico):")
	fmt.Println(`
type UserRepository interface {
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByPhone(ctx context.Context, phone string) (*User, error)
    FindByRole(ctx context.Context, role Role) ([]*User, error)
    FindByStatus(ctx context.Context, status Status) ([]*User, error)
    FindActiveUsers(ctx context.Context) ([]*User, error)
    FindInactiveUsers(ctx context.Context) ([]*User, error)
    FindByCreatedAtBetween(ctx context.Context, start, end time.Time) ([]*User, error)
    FindByRoleIn(ctx context.Context, roles []Role) ([]*User, error)
    FindByNameLike(ctx context.Context, name string) ([]*User, error)
    // ... 50+ métodos específicos
}`)

	fmt.Println("\n✅ AGORA (Repository genérico + QueryBuilder + Specification):")
	fmt.Println(`
type Repository[T any] interface {
    Create(ctx context.Context, entity T) error
    FindOne(ctx context.Context, filter QueryFilter) (T, error)
    FindMany(ctx context.Context, filter QueryFilter) ([]T, error)
    Update(ctx context.Context, id uuid.UUID, entity T) error
    Delete(ctx context.Context, id uuid.UUID) error
    Paginate(ctx context.Context, filter QueryFilter) (*PaginatedResult[T], error)
    Count(ctx context.Context, filter QueryFilter) (int64, error)
    Exists(ctx context.Context, filter QueryFilter) (bool, error)
    // ... poucos métodos poderosos
}

// QueryBuilder - Fácil e legível
filter := NewQueryBuilder().
    WhereEqual("status", "active").
    WhereIn("role", []interface{}{"admin", "user"}).
    OrderByDesc("created_at").
    Page(1).
    PageSize(20).
    Build()

// Specification - Reutilizável
activeAdmins := ActiveAdminsSpecification[User]()
users, err := repo.FindMany(ctx, activeAdmins.ToQueryFilter())
`)

	fmt.Println("\n🎯 VANTAGENS DO PADRÃO ATUAL:")
	fmt.Println("✅ Interface enxuta (10 métodos vs 50+)")
	fmt.Println("✅ 100% flexível (qualquer busca possível)")
	fmt.Println("✅ Fácil de manter")
	fmt.Println("✅ Padrão da indústria")
	fmt.Println("✅ Type-safe com Go Generics")
	fmt.Println("✅ Reutilizável para TODAS entidades")
	fmt.Println("✅ Paginação profissional inclusa")
	fmt.Println("✅ Agregações poderosas")
	fmt.Println("✅ Transações simples")
	fmt.Println("✅ QueryBuilder facilita uso")
	fmt.Println("✅ Specification reutiliza regras")
	fmt.Println("✅ Performance otimizada com índices")
	fmt.Println("✅ Conversões seguras Domain ↔ Model")
	fmt.Println("✅ Hooks GORM para validações")
	fmt.Println("✅ Queries específicas otimizadas")
}

func main() {
	ExemploUsoRepository()
	ExemploComparacaoAntesDepois()
}
