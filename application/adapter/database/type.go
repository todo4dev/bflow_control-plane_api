package database

import (
	"context"
	"encoding/json"

	"src/core/builder"
	"src/core/common"
)

type IDatabaseAdapter interface {
	Ping(
		ctx context.Context,
	) error

	BeginTransaction(
		ctx context.Context,
	) (common.IUnitOfWork, error)

	FindOne(
		ctx context.Context,
		table string,
		query *builder.Query[json.RawMessage],
		optionalUow ...common.IUnitOfWork,
	) (*json.RawMessage, error)

	FindMany(
		ctx context.Context,
		table string,
		query *builder.Query[json.RawMessage],
		optionalUow ...common.IUnitOfWork,
	) (*builder.Result[json.RawMessage], error)

	Count(
		ctx context.Context,
		table string,
		query *builder.Query[json.RawMessage],
		optionalUow ...common.IUnitOfWork,
	) (int64, error)

	Insert(
		ctx context.Context,
		table string,
		entities []json.RawMessage,
		optionalUow ...common.IUnitOfWork,
	) error

	Update(
		ctx context.Context,
		table string,
		where *builder.WhereBuilder[json.RawMessage],
		update *builder.UpdateBuilder[json.RawMessage],
		optionalUow ...common.IUnitOfWork,
	) (int64, error)

	Delete(
		ctx context.Context,
		table string,
		where *builder.WhereBuilder[json.RawMessage],
		optionalUow ...common.IUnitOfWork,
	) (int64, error)
}
