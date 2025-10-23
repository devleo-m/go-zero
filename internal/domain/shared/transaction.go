// internal/domain/shared/transaction.go
package shared

import "context"

// TransactionManager gerencia transações
type TransactionManager interface {
	// BeginTransaction inicia uma nova transação
	BeginTransaction(ctx context.Context) (context.Context, error)

	// Commit confirma a transação
	Commit(ctx context.Context) error

	// Rollback desfaz a transação
	Rollback(ctx context.Context) error

	// WithTransaction executa função dentro de transação
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// TransactionKey é a chave usada para armazenar transação no context
type TransactionKey struct{}

// GetTransactionFromContext extrai transação do context
func GetTransactionFromContext(ctx context.Context) interface{} {
	return ctx.Value(TransactionKey{})
}

// WithTransactionInContext adiciona transação ao context
func WithTransactionInContext(ctx context.Context, tx interface{}) context.Context {
	return context.WithValue(ctx, TransactionKey{}, tx)
}
