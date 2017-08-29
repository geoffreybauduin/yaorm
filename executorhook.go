package yaorm

// ExecutorHook can be implemented and set on the database handler
// Before every query, a new object will be created to provide a way to perform operations,
// like logging / timing your sql queries
type ExecutorHook interface {
	BeforeSelectOne(query string, args ...interface{})
	AfterSelectOne(query string, args ...interface{})
	BeforeSelect(query string, args ...interface{})
	AfterSelect(query string, args ...interface{})
}

// DefaultExecutorHook is the default executor hook, returned if no hook has been defined
// This struct can be composed by any of your executor hooks to avoid having to define every handler
type DefaultExecutorHook struct{}

func (h DefaultExecutorHook) BeforeSelectOne(query string, args ...interface{}) {}
func (h DefaultExecutorHook) AfterSelectOne(query string, args ...interface{})  {}

func (h DefaultExecutorHook) BeforeSelect(query string, args ...interface{}) {}
func (h DefaultExecutorHook) AfterSelect(query string, args ...interface{})  {}
