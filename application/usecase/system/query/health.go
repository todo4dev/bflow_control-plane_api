package query

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"golang.org/x/sync/errgroup"

	"src/application/adapter/cache"
	"src/application/adapter/database"
	"src/application/adapter/storage"
	"src/application/adapter/stream"
	"src/core"
	"src/core/cqrs"
	"src/core/doc"
	"src/domain/exception"
)

// #region error

const (
	Err_HealthQuery_Failed = "health query failed"
)

// #endregion
// #region HealthQuery

type HealthQuery struct {
}

// #endregion
// #region HealthQueryResult

type HealthQueryResult struct {
	Database *string `json:"database"`
	Cache    *string `json:"cache"`
	Storage  *string `json:"storage"`
	Stream   *string `json:"stream"`
}

// #endregion
// #region HealthQueryHandler

type HealthQueryHandler struct {
	database database.IDatabaseAdapter
	cache    cache.ICacheAdapter
	storage  storage.IStorageAdapter
	stream   stream.IStreamAdapter
}

var _ cqrs.IQueryHandler[
	*HealthQuery,
	*HealthQueryResult,
] = (*HealthQueryHandler)(nil)

func NewHealthQueryHandler(
	database database.IDatabaseAdapter,
	cache cache.ICacheAdapter,
	storage storage.IStorageAdapter,
	stream stream.IStreamAdapter,
) *HealthQueryHandler {
	return &HealthQueryHandler{
		database: database,
		cache:    cache,
		storage:  storage,
		stream:   stream,
	}
}

func (h *HealthQueryHandler) pingAndSetTime(
	ctx context.Context,
	g *errgroup.Group,
	p interface {
		Ping(ctx context.Context) error
	},
	target **string,
) {
	g.Go(func() (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				*target = nil
				err = fmt.Errorf("panic in Ping: %v\n%s", recovered, debug.Stack())
			}
		}()

		start := time.Now()

		if err := p.Ping(ctx); err != nil {
			*target = nil
			return err
		}
		duration := fmt.Sprintf("%.3fms", float64(time.Since(start))/float64(time.Millisecond))
		*target = &duration
		return nil
	})
}

func (h *HealthQueryHandler) Handle(
	ctx context.Context,
	query *HealthQuery,
) (*HealthQueryResult, error) {
	var result HealthQueryResult

	g, ctx := errgroup.WithContext(ctx)
	h.pingAndSetTime(ctx, g, h.database, &result.Database)
	h.pingAndSetTime(ctx, g, h.cache, &result.Cache)
	h.pingAndSetTime(ctx, g, h.storage, &result.Storage)
	h.pingAndSetTime(ctx, g, h.stream, &result.Stream)

	if err := g.Wait(); err != nil {
		return nil, exception.NewInternalException().
			WithCause(err).
			WithMessage(Err_HealthQuery_Failed)
	}

	return &result, nil
}

// #endregion

func init() {
	query := HealthQuery{}
	doc.Describe(&query,
		doc.Description("Health check query"),
		doc.Throws[exception.InternalException](Err_HealthQuery_Failed))

	result := HealthQueryResult{
		Database: core.Ptr("123ms"),
		Cache:    core.Ptr("123ms"),
		Storage:  core.Ptr("123ms"),
		Stream:   core.Ptr("123ms"),
	}
	doc.Describe(&result,
		doc.Description("Health check result"),
		doc.Example(&result),
		doc.Field(&result.Database, doc.Description("Database connection status and response time in milliseconds, nil if failed")),
		doc.Field(&result.Cache, doc.Description("Cache connection status and response time in milliseconds, nil if failed")),
		doc.Field(&result.Storage, doc.Description("Storage connection status and response time in milliseconds, nil if failed")),
		doc.Field(&result.Stream, doc.Description("Stream connection status and response time in milliseconds, nil if failed")))

	cqrs.RegisterQueryHandler[
		*HealthQuery,
		*HealthQueryResult,
		*HealthQueryHandler,
	](NewHealthQueryHandler)
}
