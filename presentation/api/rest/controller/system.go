package controller

import (
	"net/http"

	"src/application/usecase/system/query"
	"src/core/cqrs"
	"src/core/di"
	"src/core/meta"
	"src/presentation/api/rest/core"
	"src/presentation/api/rest/interceptor"
	"src/presentation/api/rest/oas"
)

type SystemController struct {
	tags string
}

var _ core.IRestController = (*SystemController)(nil)

func NewSystemController() *SystemController {
	return &SystemController{tags: "System"}
}

func (c *SystemController) Router() core.Router {
	return core.NewRouter().
		Push(c.GetHealth())
}

func (c *SystemController) GetHealth() *core.RouteBuilder {
	metadata := meta.GetObjectMetadataAs[query.HealthcheckQuery]()
	return core.NewRoute().Get("/health").
		OperationId("SystemHealth").Tags(c.tags).
		Summary(metadata.Description).Description(metadata.Description).
		Response(http.StatusOK, func(r *oas.BuildResponse) {
			metadata := meta.GetObjectMetadataAs[query.HealthcheckQueryResult]()
			r.Description(metadata.Description).Content(oas.ContentType_ApplicationJson, func(m *oas.BuildMediaType) {
				m.Schema(oas.ObjectMetadata(metadata)).Example(metadata.Example)
			})
		}).
		Handler(func(ctx core.HttpContext) error {
			result := cqrs.MustExecuteQuery[query.HealthcheckQueryResult](ctx.Context(), &query.HealthcheckQuery{})
			return ctx.JSON(http.StatusOK, result)
		}).
		UseInterceptors(interceptor.LoggingInterceptor())
}

func init() {
	di.RegisterAs[core.IRestController](NewSystemController)
}
