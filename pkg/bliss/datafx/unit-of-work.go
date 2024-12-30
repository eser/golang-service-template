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

type UnitOfWork interface {
	TxScope() TransactionFinalizer
	Context() context.Context
	Bind(context context.Context, txScope TransactionFinalizer)
	Commit() error
	Close() error
}

type UnitOfWorkImpl struct {
	context context.Context //nolint:containedctx
	txScope TransactionFinalizer
}

var _ UnitOfWork = (*UnitOfWorkImpl)(nil)

// func NewUnitOfWork() *UnitOfWorkImpl {
// 	return &UnitOfWorkImpl{}
// }

type TransactionStarter interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func UseUnitOfWork(ctx context.Context, transactionStarter TransactionStarter) (UnitOfWork, error) { //nolint:ireturn,varnamelen,lll
	uow, ok := ctx.Value(ContextKeyUnitOfWork).(UnitOfWork)
	if ok {
		return uow, nil
	}

	uow = &UnitOfWorkImpl{} //nolint:exhaustruct
	newCtx := context.WithValue(ctx, ContextKeyUnitOfWork, uow)

	transaction, err := transactionStarter.BeginTx(newCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	uow.Bind(newCtx, transaction)

	return uow, nil
}

func (uow *UnitOfWorkImpl) TxScope() TransactionFinalizer { //nolint:ireturn
	return uow.txScope
}

func (uow *UnitOfWorkImpl) Context() context.Context {
	return uow.context
}

func (uow *UnitOfWorkImpl) Bind(context context.Context, txScope TransactionFinalizer) {
	uow.context = context
	uow.txScope = txScope
}

func (uow *UnitOfWorkImpl) Commit() error {
	return uow.txScope.Commit() //nolint:wrapcheck
}

func (uow *UnitOfWorkImpl) Close() error {
	return uow.txScope.Rollback() //nolint:wrapcheck
}

func (uow *UnitOfWorkImpl) Use(fn func(TransactionFinalizer) any) {
	fn(uow.txScope)
}
