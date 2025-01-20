package datafx

import (
	"context"
	"database/sql"
	"fmt"
)

type ContextKey string

const (
	ContextKeyUnitOfWork ContextKey = "unit-of-work"
)

type TransactionFinalizer interface {
	Rollback() error
	Commit() error
}

type UnitOfWork struct {
	context context.Context //nolint:containedctx
	txScope TransactionFinalizer
}

// func NewUnitOfWork() *UnitOfWork {
// 	return &UnitOfWork{}
// }

type TransactionStarter interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func UseUnitOfWork(ctx context.Context, transactionStarter TransactionStarter) (*UnitOfWork, error) {
	uow, ok := ctx.Value(ContextKeyUnitOfWork).(*UnitOfWork)
	if ok {
		return uow, nil
	}

	uow = &UnitOfWork{} //nolint:exhaustruct
	newCtx := context.WithValue(ctx, ContextKeyUnitOfWork, uow)

	transaction, err := transactionStarter.BeginTx(newCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	uow.Bind(newCtx, transaction)

	return uow, nil
}

func (uow *UnitOfWork) TxScope() TransactionFinalizer { //nolint:ireturn
	return uow.txScope
}

func (uow *UnitOfWork) Context() context.Context {
	return uow.context
}

func (uow *UnitOfWork) Bind(context context.Context, txScope TransactionFinalizer) {
	uow.context = context
	uow.txScope = txScope
}

func (uow *UnitOfWork) Commit() error {
	return uow.txScope.Commit() //nolint:wrapcheck
}

func (uow *UnitOfWork) Close() error {
	return uow.txScope.Rollback() //nolint:wrapcheck
}

func (uow *UnitOfWork) Use(fn func(TransactionFinalizer) any) {
	fn(uow.txScope)
}
