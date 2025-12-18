package common

import "context"

type IUnitOfWork interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
