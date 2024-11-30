package datafx

import (
	"context"
	"fmt"
)

type ContextKey string

const (
	ContextKeyUnitOfWork ContextKey = "unit-of-work"
)

type UnitOfWork interface {
	Scope() DbExecutorTx
	Context() context.Context
	Bind(context context.Context, scope DbExecutorTx)
	Commit() error
	Close() error
}

type UnitOfWorkImpl struct {
	context context.Context //nolint:containedctx
	scope   DbExecutorTx
}

var _ UnitOfWork = (*UnitOfWorkImpl)(nil)

// func NewUnitOfWork() *UnitOfWorkImpl {
// 	return &UnitOfWorkImpl{}
// }

func UseUnitOfWork(ctx context.Context, db DbExecutorDb) (UnitOfWork, error) { //nolint:ireturn,varnamelen
	uow, ok := ctx.Value(ContextKeyUnitOfWork).(UnitOfWork)
	if ok {
		return uow, nil
	}

	uow = &UnitOfWorkImpl{} //nolint:exhaustruct
	newCtx := context.WithValue(ctx, ContextKeyUnitOfWork, uow)

	transaction, err := db.BeginTx(newCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	uow.Bind(newCtx, transaction)

	return uow, nil
}

func (uow *UnitOfWorkImpl) Scope() DbExecutorTx { //nolint:ireturn
	return uow.scope
}

func (uow *UnitOfWorkImpl) Context() context.Context {
	return uow.context
}

func (uow *UnitOfWorkImpl) Bind(context context.Context, scope DbExecutorTx) {
	uow.context = context
	uow.scope = scope
}

func (uow *UnitOfWorkImpl) Commit() error {
	return uow.scope.Commit() //nolint:wrapcheck
}

func (uow *UnitOfWorkImpl) Close() error {
	return uow.scope.Rollback() //nolint:wrapcheck
}

func (uow *UnitOfWorkImpl) Use(fn func(DbExecutor) any) {
	fn(uow.scope)
}
