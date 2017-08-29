package yaormfilter

// NilFilter is a filter to operate on unknown nil fields
type NilFilter struct {
	valuefilterimpl
}

// NewNilFilter returns a new NilFilter
func NewNilFilter() ValueFilter {
	return &NilFilter{}
}

// Equals applies an equal filter on Date
func (f *NilFilter) Equals(v interface{}) ValueFilter {
	return f
}

// Like is not applicable on Date
func (f *NilFilter) Like(v interface{}) ValueFilter {
	return f
}

// Nil applies a nil filter on Date
func (f *NilFilter) Nil(v bool) ValueFilter {
	f.nil(v)
	return f

}

// In adds a IN filter (not implemented)
func (f *NilFilter) In(values ...interface{}) ValueFilter {
	return f
}
