package api

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"src/application/adapter/logger"
	"src/core/di"
	"src/core/env"
	"src/presentation/api/rest/core"
	"src/presentation/api/rest/oas"

	_ "src/presentation/api/rest/controller"
	_ "src/presentation/api/rest/interceptor"
)

type Server struct {
	Config *core.Config
	app    *fiber.App
	routes []core.Route
}

func NewServer(config *core.Config) *Server {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	server := &Server{Config: config}

	server.app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err == nil {
				return nil
			}

			ctx := core.NewFiberHttpContext(c)
			logger := di.Resolve[logger.ILoggerAdapter]()
			logger.Error(err.Error(), map[string]any{
				"method": c.Method(),
				"path":   c.Path(),
			})

			status := core.GetHTTPStatus(err)

			return ctx.JSON(status, map[string]any{
				"error": err.Error(),
			})
		},
	})

	server.registerOpenAPI()
	server.loadControllers()

	if config.EnableSwagger {
		server.enableSwaggerUIHandler()
		server.enableSwaggerJSONHandler()
	}

	return server
}

func (s *Server) registerOpenAPI() {
	di.SingletonAs[oas.OpenAPI](func() oas.OpenAPI {
		return oas.NewOpenAPI(func(b *oas.BuildOpenAPI) {
			b.Info(s.Config.SwaggerTitle, s.Config.SwaggerVersion, func(i *oas.BuildInfo) {
				i.Description(s.Config.SwaggerDescription).
					Contact(s.Config.SwaggerContactName, s.Config.SwaggerContactURL, s.Config.SwaggerContactEmail).
					License(s.Config.SwaggerLicenseName, s.Config.SwaggerLicenseURL)
			})
			b.Server("http://localhost:"+strconv.Itoa(s.Config.Port), "Localhost")
		})
	})
}

func (s *Server) chainInterceptors(final core.HandlerFN, interceptors ...core.InterceptorFN) core.HandlerFN {
	if len(interceptors) == 0 {
		return final
	}

	handler := final
	for i := len(interceptors) - 1; i >= 0; i-- {
		current := interceptors[i]
		next := handler

		handler = func(ctx core.HttpContext) error {
			return current(ctx, next)
		}
	}

	return handler
}

func (s *Server) loadControllers() {
	controllers := di.ResolveAll[core.IRestController]()
	var routes []core.Route
	for _, controller := range controllers {
		router := controller.Router()
		for _, r := range router.Routes {
			route := r
			fullPath := path.Clean(s.Config.BasePath + router.BasePath + route.Path)
			route.Path = fullPath
			routes = append(routes, route)
			s.registerRoute(route)
		}
	}

	s.routes = routes
}

func (s *Server) registerRoute(route core.Route) {
	handler := func(c *fiber.Ctx) error {
		ctx := core.NewFiberHttpContext(c)

		for _, guard := range route.Guards {
			if err := guard(ctx); err != nil {
				return err // cai no ErrorHandler global
			}
		}

		finalHandler := route.Handler
		if len(route.Interceptors) > 0 {
			finalHandler = s.chainInterceptors(finalHandler, route.Interceptors...)
		}

		return finalHandler(ctx)
	}

	switch route.Method {
	case core.Route_Get:
		s.app.Get(route.Path, handler)
	case core.Route_Post:
		s.app.Post(route.Path, handler)
	case core.Route_Put:
		s.app.Put(route.Path, handler)
	case core.Route_Delete:
		s.app.Delete(route.Path, handler)
	case core.Route_Patch:
		s.app.Patch(route.Path, handler)
	case core.Route_Options:
		s.app.Options(route.Path, handler)
	case core.Route_Head:
		s.app.Head(route.Path, handler)
	default:
		panic("method not supported")
	}
}

func (s *Server) enableSwaggerUIHandler() {
	swaggerUIPath := path.Clean(s.Config.BasePath + s.Config.SwaggerPath)
	swaggerJSONPath := path.Clean(swaggerUIPath + "/openapi.json")

	s.app.Get(swaggerUIPath, func(c *fiber.Ctx) error {
		html := fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
				<head>
					<meta charset="utf-8">
					<title>%s</title>
					<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
				</head>
				<body>
					<div id="swagger-ui"></div>
					<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
					<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-standalone-preset.js"></script>
					<script>
						window.onload = function() {
							window.ui = SwaggerUIBundle({
								url: %q,
								dom_id: '#swagger-ui',
								presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
								layout: "BaseLayout",
							});
						};
					</script>
				</body>
			</html>`,
			s.Config.SwaggerTitle,
			swaggerJSONPath,
		)

		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(html)
	})
}

func (s *Server) enableSwaggerJSONHandler() {
	swaggerUIPath := path.Clean(s.Config.BasePath + s.Config.SwaggerPath)
	swaggerJSONPath := path.Clean(swaggerUIPath + "/openapi.json")

	s.app.Get(swaggerJSONPath, func(c *fiber.Ctx) error {
		return c.JSON(di.Resolve[oas.OpenAPI]())
	})
}

func (s *Server) Routes() []core.Route {
	return s.routes
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	errChannel := make(chan error, 1)

	go func() {
		errChannel <- s.app.Listen(":" + strconv.Itoa(s.Config.Port))
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = s.app.Shutdown()
		<-shutdownCtx.Done()

		return ctx.Err()
	case err := <-errChannel:
		return err
	}
}

func init() {
	di.Register(func() *Server {
		config := &core.Config{
			Port:          env.Get("SERVER_PORT", 4000),
			BasePath:      env.Get("SERVER_BASE_PATH", "/"),
			EnableSwagger: env.Get("SERVER_ENABLE_SWAGGER", "true") == "true",
			SwaggerPath:   env.Get("SERVER_SWAGGER_PATH", "/swagger"),
		}
		return NewServer(config)
	})
}
