package yaorm

import "context"

// ExecutorHook can be implemented and set on the database handler
// Before every query, a new object will be created to provide a way to perform operations,
// like logging / timing your sql queries
type ExecutorHook interface {
	BeforeSelectOne(ctx context.Context, query string, args ...interface{})
	AfterSelectOne(ctx context.Context, query string, args ...interface{})
	BeforeSelect(ctx context.Context, query string, args ...interface{})
	AfterSelect(ctx context.Context, query string, args ...interface{})
}

// DefaultExecutorHook is the default executor hook, returned if no hook has been defined
// This struct can be composed by any of your executor hooks to avoid having to define every handler
type DefaultExecutorHook struct{}

func (h DefaultExecutorHook) BeforeSelectOne(ctx context.Context, query string, args ...interface{}) {}
func (h DefaultExecutorHook) AfterSelectOne(ctx context.Context, query string, args ...interface{})  {}

func (h DefaultExecutorHook) BeforeSelect(ctx context.Context, query string, args ...interface{}) {}
func (h DefaultExecutorHook) AfterSelect(ctx context.Context, query string, args ...interface{})  {}
