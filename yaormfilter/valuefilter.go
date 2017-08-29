package yaormfilter

import (
	"fmt"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
)

type ValueFilter interface {
	Apply(statement squirrel.SelectBuilder, tableName, fieldName string) squirrel.SelectBuilder
	Equals(v interface{}) ValueFilter
	Like(v interface{}) ValueFilter
	Nil(v bool) ValueFilter
	IsEquality() bool
	GetEquality() interface{}
	In(v ...interface{}) ValueFilter
}

type valuefilterimpl struct {
	equals_     interface{}
	nil_        *bool
	shouldEqual bool
	like_       interface{}
	shouldLike  bool
	in_         []interface{}
	shouldIn    bool
}

func (f valuefilterimpl) IsEquality() bool {
	return f.shouldEqual
}

func (f valuefilterimpl) GetEquality() interface{} {
	return f.equals_
}

func (f *valuefilterimpl) nil(v bool) *valuefilterimpl {
	f.nil_ = &v
	return f
}

func (f *valuefilterimpl) equals(e interface{}) *valuefilterimpl {
	f.equals_ = e
	f.shouldEqual = true
	return f
}

func (f *valuefilterimpl) like(e interface{}) *valuefilterimpl {
	f.like_ = e
	f.shouldLike = true
	return f
}

func (f *valuefilterimpl) in(e []interface{}) *valuefilterimpl {
	f.in_ = e
	f.shouldIn = true
	return f
}

func (f *valuefilterimpl) Apply(statement squirrel.SelectBuilder, tableName, fieldName string) squirrel.SelectBuilder {
	computedField := fmt.Sprintf(`%s.%s`, tableName, fieldName)
	if f.nil_ != nil {
		if *f.nil_ == true {
			statement = statement.Where(
				squirrel.Eq{computedField: nil},
			)
		} else {
			statement = statement.Where(
				squirrel.NotEq{computedField: nil},
			)
		}
	}
	if f.IsEquality() {
		statement = statement.Where(
			squirrel.Eq{computedField: f.GetEquality()},
		)
	}
	if f.shouldLike {
		statement = statement.Where(
			fmt.Sprintf("%s LIKE ?", computedField), f.like_,
		)
	}
	if f.shouldIn {
		statement = statement.Where(
			squirrel.Eq{computedField: f.in_},
		)
	}
	return statement
}
