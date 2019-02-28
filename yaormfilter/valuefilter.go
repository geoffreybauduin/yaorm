package yaormfilter

import (
	"fmt"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
)

type ValueFilter interface {
	Apply(statement squirrel.SelectBuilder, tableName, fieldName string) squirrel.SelectBuilder
	Equals(v interface{}) ValueFilter
	NotEquals(v interface{}) ValueFilter
	Like(v interface{}) ValueFilter
	Lt(v interface{}) ValueFilter
	Lte(v interface{}) ValueFilter
	Gt(v interface{}) ValueFilter
	Gte(v interface{}) ValueFilter
	Nil(v bool) ValueFilter
	IsEquality() bool
	GetEquality() interface{}
	In(v ...interface{}) ValueFilter
}

type valuefilterimpl struct {
	equals_        interface{}
	nil_           *bool
	shouldEqual    bool
	like_          interface{}
	shouldLike     bool
	in_            []interface{}
	shouldIn       bool
	lt_            interface{}
	shouldLt       bool
	lte_           interface{}
	shouldLte      bool
	gt_            interface{}
	shouldGt       bool
	gte_           interface{}
	shouldGte      bool
	notEquals_     interface{}
	shouldNotEqual bool
}

func (f valuefilterimpl) IsEquality() bool {
	return f.shouldEqual
}

func (f valuefilterimpl) IsInequality() bool {
	return f.shouldNotEqual
}

func (f valuefilterimpl) GetEquality() interface{} {
	return f.equals_
}

func (f valuefilterimpl) GetInequality() interface{} {
	return f.notEquals_
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

func (f *valuefilterimpl) notEquals(e interface{}) *valuefilterimpl {
	f.notEquals_ = e
	f.shouldNotEqual = true
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

func (f *valuefilterimpl) lte(e interface{}) *valuefilterimpl {
	f.lte_ = e
	f.shouldLte = true
	return f
}

func (f *valuefilterimpl) gte(e interface{}) *valuefilterimpl {
	f.gte_ = e
	f.shouldGte = true
	return f
}

func (f *valuefilterimpl) lt(e interface{}) *valuefilterimpl {
	f.lt_ = e
	f.shouldLt = true
	return f
}

func (f *valuefilterimpl) gt(e interface{}) *valuefilterimpl {
	f.gt_ = e
	f.shouldGt = true
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
	if f.IsInequality() {
		statement = statement.Where(
			squirrel.NotEq{computedField: f.GetInequality()},
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
	if f.shouldLt {
		statement = statement.Where(
			squirrel.Lt{computedField: f.lt_},
		)
	}
	if f.shouldLte {
		statement = statement.Where(
			squirrel.LtOrEq{computedField: f.lte_},
		)
	}
	if f.shouldGt {
		statement = statement.Where(
			squirrel.Gt{computedField: f.gt_},
		)
	}
	if f.shouldGte {
		statement = statement.Where(
			squirrel.GtOrEq{computedField: f.gte_},
		)
	}
	return statement
}
