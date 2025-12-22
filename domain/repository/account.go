package repository

import (
	"context"

	"src/core/common"
	"src/domain/entity"
)

type IAccountRepository interface {
	CountByEmail(ctx context.Context, email string, optionalUow ...common.IUnitOfWork) (int64, error)
	GetByEmail(ctx context.Context, email string, optionalUow ...common.IUnitOfWork) (*entity.AccountEntity, error)
	ActivateByEmail(ctx context.Context, email string, optionalUow ...common.IUnitOfWork) (*entity.AccountEntity, error)
}
