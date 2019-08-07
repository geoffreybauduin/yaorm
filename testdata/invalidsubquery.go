package testdata

import (
	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type InvalidSubquery struct {
	yaorm.DatabaseModel
	ID int64 `db:"id"`
}

type InvalidSubqueryFilter struct {
	yaormfilter.ModelFilter
	FilterID yaormfilter.ValueFilter `filter:"id"`
}

func init() {
	yaorm.NewTable("test", "invalid_subquery", &InvalidSubquery{}).WithFilter(&InvalidSubqueryFilter{}).WithSubqueryloading(
		func(dbp yaorm.DBProvider, ids []interface{}) (interface{}, error) {
			return yaorm.GenericSelectOne(dbp, NewInvalidSubqueryFilter().ID(yaormfilter.In(ids...)))
		},
		"id",
	)
	yaorm.NewTable("test", "invalid_subquery_b", &InvalidSubqueryB{}).WithFilter(&InvalidSubqueryBFilter{})
}

func NewInvalidSubqueryFilter() *InvalidSubqueryFilter {
	return &InvalidSubqueryFilter{}
}

func (f *InvalidSubqueryFilter) ID(v yaormfilter.ValueFilter) *InvalidSubqueryFilter {
	f.FilterID = v
	return f
}

func (f *InvalidSubqueryFilter) Subqueryload() yaormfilter.Filter {
	f.AllowSubqueryload()
	return f
}

type InvalidSubqueryB struct {
	yaorm.DatabaseModel
	ID                int64            `db:"id"`
	InvalidSubqueryID int64            `db:"invalid_subquery_id"`
	InvalidSubquery   *InvalidSubquery `db:"-" filterload:"invalid_subquery,invalid_subquery_id"`
}

type InvalidSubqueryBFilter struct {
	yaormfilter.ModelFilter
	FilterInvalidSubquery yaormfilter.Filter `filter:"invalid_subquery,join,id,invalid_subquery_id" filterload:"invalid_subquery"`
}

func NewInvalidSubqueryBFilter() *InvalidSubqueryBFilter {
	return &InvalidSubqueryBFilter{}
}

func (f *InvalidSubqueryBFilter) InvalidSubquery(v yaormfilter.Filter) *InvalidSubqueryBFilter {
	if _, ok := v.(*InvalidSubqueryFilter); !ok {
		panic("Not a InvalidSubqueryFilter")
	}
	f.FilterInvalidSubquery = v
	return f
}
