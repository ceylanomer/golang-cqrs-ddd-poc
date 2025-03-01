package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormZapLogger struct {
	logger        *zap.Logger
	slowThreshold time.Duration
}

func NewGormZapLogger(zapLogger *zap.Logger) gormlogger.Interface {
	return &GormZapLogger{
		logger:        zapLogger,
		slowThreshold: 200 * time.Millisecond,
	}
}

func (l *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, data...))
}

func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, data...))
}

func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, data...))
}

func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fields = append(fields, zap.Error(err))
		l.logger.Error("gorm-trace", fields...)
		return
	}

	if elapsed > l.slowThreshold {
		l.logger.Warn("gorm-trace-slow-query", fields...)
		return
	}

	l.logger.Debug("gorm-trace", fields...)
}
