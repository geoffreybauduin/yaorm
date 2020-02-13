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
	BeforeInsert(ctx context.Context, query string, args ...interface{})
	AfterInsert(ctx context.Context, query string, args ...interface{})
	BeforeUpdate(ctx context.Context, query string, args ...interface{})
	AfterUpdate(ctx context.Context, query string, args ...interface{})
	BeforeDelete(ctx context.Context, query string, args ...interface{})
	AfterDelete(ctx context.Context, query string, args ...interface{})
	BeforeExec(ctx context.Context, query string, args ...interface{})
	AfterExec(ctx context.Context, query string, args ...interface{})
}

// DefaultExecutorHook is the default executor hook, returned if no hook has been defined
// This struct can be composed by any of your executor hooks to avoid having to define every handler
type DefaultExecutorHook struct{}

// BeforeSelectOne is a hook executed before performing a SelectOne action on gorp sql executor
func (h DefaultExecutorHook) BeforeSelectOne(ctx context.Context, query string, args ...interface{}) {}

// AfterSelectOne is a hook executed before performing a SelectOne action on gorp sql executor
func (h DefaultExecutorHook) AfterSelectOne(ctx context.Context, query string, args ...interface{}) {}

// BeforeSelect is a hook executed before performing a Select action on gorp sql executor
func (h DefaultExecutorHook) BeforeSelect(ctx context.Context, query string, args ...interface{}) {}

// AfterSelect is a hook executed before performing a Select action on gorp sql executor
func (h DefaultExecutorHook) AfterSelect(ctx context.Context, query string, args ...interface{}) {}

// BeforeInsert is a hook executed before performing an Insert action on gorp sql executor
func (h DefaultExecutorHook) BeforeInsert(ctx context.Context, query string, args ...interface{}) {}

// AfterInsert is a hook executed before performing an Insert action on gorp sql executor
func (h DefaultExecutorHook) AfterInsert(ctx context.Context, query string, args ...interface{}) {}

// BeforeUpdate is a hook executed before performing an Update action on gorp sql executor
func (h DefaultExecutorHook) BeforeUpdate(ctx context.Context, query string, args ...interface{}) {}

// AfterUpdate is a hook executed before performing an Update action on gorp sql executor
func (h DefaultExecutorHook) AfterUpdate(ctx context.Context, query string, args ...interface{}) {}

// BeforeDelete is a hook executed before performing a Delete action on gorp sql executor
func (h DefaultExecutorHook) BeforeDelete(ctx context.Context, query string, args ...interface{}) {}

// AfterDelete is a hook executed before performing a Delete action on gorp sql executor
func (h DefaultExecutorHook) AfterDelete(ctx context.Context, query string, args ...interface{}) {}

// BeforeExec is a hook executed before performing an Exec action on gorp sql executor
func (h DefaultExecutorHook) BeforeExec(ctx context.Context, query string, args ...interface{}) {}

// AfterExec is a hook executed before performing an Exec action on gorp sql executor
func (h DefaultExecutorHook) AfterExec(ctx context.Context, query string, args ...interface{}) {}
