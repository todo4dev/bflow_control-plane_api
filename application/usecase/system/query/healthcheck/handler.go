package healthcheck

import (
	"context"
	"fmt"
	"src/application/adapter/cache"
	"src/application/adapter/database"
	"src/application/adapter/storage"
	"src/application/adapter/stream"
	"src/core/cqrs"
	"src/domain/exception"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	Err_Failed = "health query failed"
)

type Handler struct {
	database database.IDatabaseAdapter
	cache    cache.ICacheAdapter
	storage  storage.IStorageAdapter
	stream   stream.IStreamAdapter
}

var _ cqrs.IQueryHandler[*Query, *Result] = (*Handler)(nil)

func New(
	database database.IDatabaseAdapter,
	cache cache.ICacheAdapter,
	storage storage.IStorageAdapter,
	stream stream.IStreamAdapter,
) *Handler {
	return &Handler{
		database: database,
		cache:    cache,
		storage:  storage,
		stream:   stream,
	}
}

func (h *Handler) pingAndSetTime(
	ctx context.Context,
	g *errgroup.Group,
	p interface {
		Ping(ctx context.Context) error
	},
	target **string,
) {
	g.Go(func() error {
		start := time.Now()
		err := p.Ping(ctx)
		if err != nil {
			*target = nil
			return err
		}
		duration := fmt.Sprintf("%dms", time.Since(start).Milliseconds())
		*target = &duration
		return nil
	})
}

func (h *Handler) Handle(
	ctx context.Context,
	query *Query,
) (*Result, error) {
	var result Result
	g, ctx := errgroup.WithContext(ctx)
	h.pingAndSetTime(ctx, g, h.database, &result.Database)
	h.pingAndSetTime(ctx, g, h.cache, &result.Cache)
	h.pingAndSetTime(ctx, g, h.storage, &result.Storage)
	h.pingAndSetTime(ctx, g, h.stream, &result.Stream)

	if err := g.Wait(); err != nil {
		return nil, exception.NewInternal().WithCause(err).WithMessage(Err_Failed)
	}

	return &result, nil
}
