package handler

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type Request any
type Response any

// Define an interface for handlers
type HandlerInterface[R Request, Res Response] interface {
	Handle(ctx context.Context, req *R) (*Res, error)
}

// Update handle function to accept HandlerInterface instead of Handler function
func Handler[R Request, Res Response](handler HandlerInterface[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		// Start a new span for this handler
		ctx := c.UserContext()
		tracer := otel.GetTracerProvider().Tracer("")
		spanName := c.Route().Path
		ctx, span := tracer.Start(ctx, spanName)
		defer span.End()

		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			span.RecordError(err)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := c.ParamsParser(&req); err != nil {
			span.RecordError(err)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := c.QueryParser(&req); err != nil {
			span.RecordError(err)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := c.ReqHeaderParser(&req); err != nil {
			span.RecordError(err)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		/*
			ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
			defer cancel()
		*/

		res, err := handler.Handle(ctx, &req)
		if err != nil {
			if fiberErr, ok := err.(*fiber.Error); ok {
				span.RecordError(err)
				zap.L().Error("Failed to handle request", logger.GetTraceFieldsWithError(ctx, err)...)
				return fiberErr
			}
			span.RecordError(err)
			zap.L().Error("Failed to handle request", logger.GetTraceFieldsWithError(ctx, err)...)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(res)
	}
}
