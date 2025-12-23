package repository

import (
	"context"
	"encoding/json"
	"src/application/adapter/database"
	"src/core/builder"
	"src/core/common"
	"src/core/di"
	"src/domain/entity"
	"src/domain/repository"
)

type PgxAccountRepository struct {
	tableName       string
	entityType      entity.AccountEntity
	databaseAdapter database.IDatabaseAdapter
}

var _ repository.IAccountRepository = (*PgxAccountRepository)(nil)

func NewPgxAccountRepository(
	databaseAdapter database.IDatabaseAdapter,
) *PgxAccountRepository {
	return &PgxAccountRepository{
		tableName:       `"control_plane"."account"`,
		entityType:      entity.AccountEntity{},
		databaseAdapter: databaseAdapter,
	}
}

func (r *PgxAccountRepository) CountByEmail(
	ctx context.Context,
	email string,
	optionalUow ...common.IUnitOfWork,
) (int64, error) {
	return r.databaseAdapter.Count(ctx, r.tableName,
		builder.NewQuery[entity.AccountEntity]().
			Where(func(e *entity.AccountEntity, q *builder.WhereBuilder[entity.AccountEntity]) {
				q.Equal(&r.entityType.Email, email)
			}).
			ToJSON(),
		optionalUow...,
	)
}

func (r *PgxAccountRepository) GetByEmail(
	ctx context.Context,
	email string,
	optionalUow ...common.IUnitOfWork,
) (*entity.AccountEntity, error) {
	query := builder.NewQuery[entity.AccountEntity]().
		Where(func(e *entity.AccountEntity, q *builder.WhereBuilder[entity.AccountEntity]) {
			q.Equal(&r.entityType.Email, email)
		}).
		ToJSON()
	jsonAccount, err := r.databaseAdapter.FindOne(ctx, r.tableName, query, optionalUow...)
	if err != nil {
		return nil, err
	}
	var typedAccount entity.AccountEntity
	if err := json.Unmarshal(*jsonAccount, &typedAccount); err != nil {
		return nil, err
	}
	return &typedAccount, nil
}

func (r *PgxAccountRepository) ActivateByEmail(
	ctx context.Context,
	email string,
	optionalUow ...common.IUnitOfWork,
) (*entity.AccountEntity, error) {
	if _, err := r.databaseAdapter.Update(ctx, r.tableName,
		builder.NewWhere[entity.AccountEntity]().
			Equal(&r.entityType.Email, email).
			ToJSON(),
		builder.NewUpdate[entity.AccountEntity]().
			Set(&r.entityType.DeletedAt, nil).
			ToJSON(),
	); err != nil {
		return nil, err
	}

	return database.TypedFromJsonWithErr[entity.AccountEntity](
		r.databaseAdapter.FindOne(ctx, r.tableName,
			builder.NewQuery[entity.AccountEntity]().
				Where(func(e *entity.AccountEntity, q *builder.WhereBuilder[entity.AccountEntity]) {
					q.Equal(&r.entityType.Email, email)
				}).
				ToJSON(),
			optionalUow...,
		),
	)
}

func init() {
	di.SingletonAs[repository.IAccountRepository](NewPgxAccountRepository)
}
