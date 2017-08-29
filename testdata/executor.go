package testdata

import (
	"log"

	"github.com/geoffreybauduin/yaorm"
)

type LoggingExecutor struct {
	yaorm.DefaultExecutorHook
}

func (l LoggingExecutor) AfterSelectOne(query string, args ...interface{}) {
	l.logQuery(query, args...)
}

func (l LoggingExecutor) AfterSelect(query string, args ...interface{}) {
	l.logQuery(query, args...)
}

func (l LoggingExecutor) logQuery(query string, args ...interface{}) {
	log.Printf("Query: %s %+v\n", query, args)
}
