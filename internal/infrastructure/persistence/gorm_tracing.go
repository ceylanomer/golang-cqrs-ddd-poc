package persistence

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	tracerName = "gorm"
)

// RegisterGormTracing registers OpenTelemetry tracing callbacks with GORM
func RegisterGormTracing(db *gorm.DB, tp trace.TracerProvider) error {
	if tp == nil {
		tp = otel.GetTracerProvider()
	}

	tracer := tp.Tracer(tracerName)

	err := db.Callback().Create().Before("gorm:create").Register("otel:before_create", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.create")
	})
	if err != nil {
		return fmt.Errorf("failed to register create callback: %w", err)
	}

	err = db.Callback().Query().Before("gorm:query").Register("otel:before_query", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.query")
	})
	if err != nil {
		return fmt.Errorf("failed to register query callback: %w", err)
	}

	err = db.Callback().Update().Before("gorm:update").Register("otel:before_update", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.update")
	})
	if err != nil {
		return fmt.Errorf("failed to register update callback: %w", err)
	}

	err = db.Callback().Delete().Before("gorm:delete").Register("otel:before_delete", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.delete")
	})
	if err != nil {
		return fmt.Errorf("failed to register delete callback: %w", err)
	}

	err = db.Callback().Row().Before("gorm:row").Register("otel:before_row", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.row")
	})
	if err != nil {
		return fmt.Errorf("failed to register row callback: %w", err)
	}

	err = db.Callback().Raw().Before("gorm:raw").Register("otel:before_raw", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.raw")
	})
	if err != nil {
		return fmt.Errorf("failed to register raw callback: %w", err)
	}

	return nil
}

func startSpan(db *gorm.DB, tracer trace.Tracer, operation string) {
	ctx := db.Statement.Context
	if ctx == nil {
		return
	}

	spanName := fmt.Sprintf("%s %s", operation, db.Statement.Table)
	opts := []trace.SpanStartOption{
		trace.WithAttributes(
			attribute.String("db.system", "postgres"),
			attribute.String("db.operation", operation),
			attribute.String("db.table", db.Statement.Table),
		),
	}

	newCtx, span := tracer.Start(ctx, spanName, opts...)
	db.Statement.Context = newCtx

	db.Callback().Create().After("*").Register("otel:after", func(db *gorm.DB) {
		defer span.End()
		if db.Error != nil {
			span.RecordError(db.Error)
		}
	})
}
