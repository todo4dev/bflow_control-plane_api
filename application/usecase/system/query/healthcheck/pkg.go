package healthcheck

import (
	"src/core"
	"src/core/cqrs"
	"src/core/meta"
	"src/domain/exception"
)

func Register() {
	registerMeta()
	cqrs.RegisterQueryHandler[*Query, *Result, *Handler](New)
}

func registerMeta() {
	query := Query{}
	meta.Describe(&query,
		meta.Description("Health check query"),
		meta.Throws[exception.Internal](Err_Failed))

	result := Result{
		Database: core.Ptr("123ms"),
		Cache:    core.Ptr("123ms"),
		Storage:  core.Ptr("123ms"),
		Stream:   core.Ptr("123ms"),
	}
	meta.Describe(&result,
		meta.Description("Health check result"),
		meta.Example(&result),
		meta.Field(&result.Database, meta.Description("Database status, null if failed")),
		meta.Field(&result.Cache, meta.Description("Cache status, null if failed")),
		meta.Field(&result.Storage, meta.Description("Storage status, null if failed")),
		meta.Field(&result.Stream, meta.Description("Stream status, null if failed")))

}
