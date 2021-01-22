package yaormfilter

import (
	"fmt"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
)

type RawFilterFunc func(string) interface{}

type ValueFilter interface {
	Apply(statement squirrel.SelectBuilder, tableName, fieldName string) squirrel.SelectBuilder
	Equals(v interface{}) ValueFilter
	NotEquals(v interface{}) ValueFilter
	Like(v interface{}) ValueFilter
	ILike(v interface{}) ValueFilter
	Lt(v interface{}) ValueFilter
	Lte(v interface{}) ValueFilter
	Gt(v interface{}) ValueFilter
	Gte(v interface{}) ValueFilter
	Nil(v bool) ValueFilter
	In(v ...interface{}) ValueFilter
	NotIn(v ...interface{}) ValueFilter
	Raw(fn RawFilterFunc) ValueFilter
	IsEquality() bool
	GetEquality() interface{}
}

type valuefilterimpl struct {
	filterFn    RawFilterFunc
	shouldEqual bool
	equals_     interface{}
}

func (f valuefilterimpl) IsEquality() bool {
	return f.shouldEqual
}

func (f valuefilterimpl) GetEquality() interface{} {
	return f.equals_
}

func (f *valuefilterimpl) nil(v bool) *valuefilterimpl {
	f.raw(func(field string) interface{} {
		if v {
			return squirrel.Eq{field: nil}
		}
		return squirrel.NotEq{field: nil}
	})
	return f
}

func (f *valuefilterimpl) equals(e interface{}) *valuefilterimpl {
	f.shouldEqual = true
	f.equals_ = e
	return f.raw(func(field string) interface{} {
		return squirrel.Eq{field: e}
	})
}

func (f *valuefilterimpl) notEquals(e interface{}) *valuefilterimpl {
	return f.raw(func(field string) interface{} {
		return squirrel.NotEq{field: e}
	})
}

func (f *valuefilterimpl) like(e interface{}) *valuefilterimpl {
	return f.raw(func(field string) interface{} {
		return squirrel.Expr(fmt.Sprintf("%s LIKE ?", field), e)
	})
}

func (f *valuefilterimpl) ilike(e interface{}) *valuefilterimpl {
	return f.raw(func(field string) interface{} {
		return squirrel.Expr(fmt.Sprintf("%s ILIKE ?", field), e)
	})
}

func (f *valuefilterimpl) in(e []interface{}) *valuefilterimpl {
	return f.raw(func(field string) interface{} {
		return squirrel.Eq{field: e}
	})
}

func (f *valuefilterimpl) notIn(e []interface{}) *valuefilterimpl {
	return f.notEquals(e)
}

func (f *valuefilterimpl) lte(e interface{}) *valuefilterimpl {
	return f.raw(func(field string) interface{} {
		return squirrel.LtOrEq{field: e}
	})
}

func (f *valuefilterimpl) gte(e interface{}) *valuefilterimpl {
	return f.raw(func(field string) interface{} {
		return squirrel.GtOrEq{field: e}
	})
}

func (f *valuefilterimpl) lt(e interface{}) *valuefilterimpl {
	return f.raw(func(field string) interface{} {
		return squirrel.Lt{field: e}
	})
}

func (f *valuefilterimpl) gt(e interface{}) *valuefilterimpl {
	return f.raw(func(field string) interface{} {
		return squirrel.Gt{field: e}
	})
}

func (f *valuefilterimpl) raw(fn RawFilterFunc) *valuefilterimpl {
	f.filterFn = fn
	return f
}

func (f *valuefilterimpl) Apply(statement squirrel.SelectBuilder, tableName, fieldName string) squirrel.SelectBuilder {
	computedField := fmt.Sprintf(`%s.%s`, tableName, fieldName)
	if f.filterFn != nil {
		statement = statement.Where(
			f.filterFn(computedField),
		)
	}
	return statement
}
