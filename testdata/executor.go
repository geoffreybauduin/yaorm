package testdata

import (
	"context"
	"log"

	"github.com/geoffreybauduin/yaorm"
)

type LoggingExecutor struct {
	yaorm.DefaultExecutorHook
}

func (l LoggingExecutor) AfterSelectOne(ctx context.Context, query string, args ...interface{}) {
	l.logQuery(ctx, query, args...)
}

func (l LoggingExecutor) AfterSelect(ctx context.Context, query string, args ...interface{}) {
	l.logQuery(ctx, query, args...)
}

func (l LoggingExecutor) logQuery(ctx context.Context, query string, args ...interface{}) {
	log.Printf("Query: %s %+v\n", query, args)
}

func (l LoggingExecutor) AfterInsert(ctx context.Context, query string, args ...interface{}) {
	l.logQuery(ctx, query, args...)
}

func (l LoggingExecutor) AfterUpdate(ctx context.Context, query string, args ...interface{}) {
	l.logQuery(ctx, query, args...)
}

func (l LoggingExecutor) AfterDelete(ctx context.Context, query string, args ...interface{}) {
	l.logQuery(ctx, query, args...)
}
