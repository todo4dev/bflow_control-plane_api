package core

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type HttpContext interface {
	Context() context.Context

	Method() string
	Path() string
	Param(name string) string
	Query(name string) string
	QueryDefault(name, defaultValue string) string
	Header(name string) string
	Body(dest any) error
	Status(code int)
	JSON(status int, body any) error
	HeaderSet(key, value string)
}

type fiberHttpContext struct {
	ctx *fiber.Ctx
}

func NewFiberHttpContext(ctx *fiber.Ctx) HttpContext {
	if ctx.UserContext() == nil {
		ctx.SetUserContext(context.Background())
	}
	return &fiberHttpContext{ctx: ctx}
}

func (c *fiberHttpContext) Context() context.Context {
	return c.ctx.UserContext()
}

func (c *fiberHttpContext) Method() string {
	return c.ctx.Method()
}

func (c *fiberHttpContext) Path() string {
	return c.ctx.Path()
}

func (c *fiberHttpContext) Param(name string) string {
	return c.ctx.Params(name)
}

func (c *fiberHttpContext) Query(name string) string {
	return c.ctx.Query(name)
}

func (c *fiberHttpContext) QueryDefault(name, def string) string {
	return c.ctx.Query(name, def)
}

func (c *fiberHttpContext) Header(name string) string {
	return c.ctx.Get(name)
}

func (c *fiberHttpContext) Body(dest any) error {
	return c.ctx.BodyParser(dest)
}

func (c *fiberHttpContext) Status(code int) {
	c.ctx.Status(code)
}

func (c *fiberHttpContext) JSON(status int, body any) error {
	return c.ctx.Status(status).JSON(body)
}

func (c *fiberHttpContext) HeaderSet(key, value string) {
	c.ctx.Set(key, value)
}
