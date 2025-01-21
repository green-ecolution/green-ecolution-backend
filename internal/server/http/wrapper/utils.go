package wrapper

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type FiberCtx struct {
	*fiber.Ctx
}

func NewFiberCtx(ctx *fiber.Ctx) *FiberCtx {
	return &FiberCtx{Ctx: ctx}
}

func (f *FiberCtx) GetLogger() *slog.Logger {
	return f.Locals("logger").(*slog.Logger)
}

func (f *FiberCtx) SetLogger(logger *slog.Logger) {
	f.Locals("logger", logger)
}

func (f *FiberCtx) WithLogger(args ...any) *slog.Logger {
	logger := f.GetLogger().With(args...)
	f.SetLogger(logger)
	return logger
}
