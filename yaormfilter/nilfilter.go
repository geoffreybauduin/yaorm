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

// NotEquals applies an notEqual filter
func (f *NilFilter) NotEquals(v interface{}) ValueFilter {
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

// NotIn adds a NOT IN filter (not implemented)
func (f *NilFilter) NotIn(values ...interface{}) ValueFilter {
	return f
}

// Lt adds a < filter
func (f *NilFilter) Lt(v interface{}) ValueFilter {
	return f
}

// Lte adds a <= filter
func (f *NilFilter) Lte(v interface{}) ValueFilter {
	return f
}

// Gt adds a > filter
func (f *NilFilter) Gt(v interface{}) ValueFilter {
	return f
}

// Gte adds a > filter
func (f *NilFilter) Gte(v interface{}) ValueFilter {
	return f
}

// Raw performs a Raw filter
func (f *NilFilter) Raw(s RawFilterFunc) ValueFilter {
	f.raw(s)
	return f
}
