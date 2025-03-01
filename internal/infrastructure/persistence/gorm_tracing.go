package persistence

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

	// Create
	if err := db.Callback().Create().Before("gorm:create").Register("otel:before_create", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.create")
	}); err != nil {
		return fmt.Errorf("failed to register before create callback: %w", err)
	}

	if err := db.Callback().Create().After("gorm:create").Register("otel:after_create", func(db *gorm.DB) {
		finishSpan(db)
	}); err != nil {
		return fmt.Errorf("failed to register after create callback: %w", err)
	}

	// Query
	if err := db.Callback().Query().Before("gorm:query").Register("otel:before_query", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.query")
	}); err != nil {
		return fmt.Errorf("failed to register before query callback: %w", err)
	}

	if err := db.Callback().Query().After("gorm:query").Register("otel:after_query", func(db *gorm.DB) {
		finishSpan(db)
	}); err != nil {
		return fmt.Errorf("failed to register after query callback: %w", err)
	}

	// Update
	if err := db.Callback().Update().Before("gorm:update").Register("otel:before_update", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.update")
	}); err != nil {
		return fmt.Errorf("failed to register before update callback: %w", err)
	}

	if err := db.Callback().Update().After("gorm:update").Register("otel:after_update", func(db *gorm.DB) {
		finishSpan(db)
	}); err != nil {
		return fmt.Errorf("failed to register after update callback: %w", err)
	}

	// Delete
	if err := db.Callback().Delete().Before("gorm:delete").Register("otel:before_delete", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.delete")
	}); err != nil {
		return fmt.Errorf("failed to register before delete callback: %w", err)
	}

	if err := db.Callback().Delete().After("gorm:delete").Register("otel:after_delete", func(db *gorm.DB) {
		finishSpan(db)
	}); err != nil {
		return fmt.Errorf("failed to register after delete callback: %w", err)
	}

	// Row
	if err := db.Callback().Row().Before("gorm:row").Register("otel:before_row", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.row")
	}); err != nil {
		return fmt.Errorf("failed to register before row callback: %w", err)
	}

	if err := db.Callback().Row().After("gorm:row").Register("otel:after_row", func(db *gorm.DB) {
		finishSpan(db)
	}); err != nil {
		return fmt.Errorf("failed to register after row callback: %w", err)
	}

	// Raw
	if err := db.Callback().Raw().Before("gorm:raw").Register("otel:before_raw", func(db *gorm.DB) {
		startSpan(db, tracer, "gorm.raw")
	}); err != nil {
		return fmt.Errorf("failed to register before raw callback: %w", err)
	}

	if err := db.Callback().Raw().After("gorm:raw").Register("otel:after_raw", func(db *gorm.DB) {
		finishSpan(db)
	}); err != nil {
		return fmt.Errorf("failed to register after raw callback: %w", err)
	}

	return nil
}

func startSpan(db *gorm.DB, tracer trace.Tracer, operation string) {
	if db.Statement.Context == nil {
		return
	}

	spanName := fmt.Sprintf("%s %s", operation, db.Statement.Table)
	opts := []trace.SpanStartOption{
		trace.WithAttributes(
			attribute.String("db.system", "postgres"),
			attribute.String("db.operation", operation),
			attribute.String("db.table", db.Statement.Table),
			attribute.String("db.statement", db.Statement.SQL.String()),
		),
		trace.WithSpanKind(trace.SpanKindClient),
	}

	_, span := tracer.Start(db.Statement.Context, spanName, opts...)
	db.Statement.Context = trace.ContextWithSpan(db.Statement.Context, span)
}

func finishSpan(db *gorm.DB) {
	if db.Statement.Context == nil {
		return
	}

	span := trace.SpanFromContext(db.Statement.Context)
	if span == nil {
		return
	}

	if db.Error != nil {
		span.SetStatus(codes.Error, db.Error.Error())
		span.RecordError(db.Error)
	} else {
		span.SetStatus(codes.Ok, "")
	}

	span.End()
}
