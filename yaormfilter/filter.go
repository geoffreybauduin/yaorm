package yaormfilter

import "github.com/juju/errors"

// Filter is the interface stating the methods to implement to correctly filter on models
type Filter interface {
	Subqueryload() Filter
	ShouldSubqueryload() bool
	AddOption(opt RequestOption) Filter
	GetSelectOptions() []RequestOption
	OrderBy(field string, way OrderingWay) Filter
	GetOrderBy() []*OrderBy
}

// OrderingWay is a custom type to have ordering
type OrderingWay string

// OrderingWays represents the Enum to have ordering
var OrderingWays = struct {
	Asc  OrderingWay
	Desc OrderingWay
}{
	Asc:  "ASC",
	Desc: "DESC",
}

// RequestOption is a custom type to have request options
type RequestOption string

// RequestOptions represents the Enum of Request options
var RequestOptions = struct {
	SelectForUpdate RequestOption
	LeftJoin        RequestOption
}{
	SelectForUpdate: "SelectForUpdate",
	LeftJoin:        "LeftJoin",
}

// ModelFilter is the struct every filter should compose
type ModelFilter struct {
	subqueryload bool
	options      []RequestOption
	orderBy      []*OrderBy
}

type OrderBy struct {
	Field string
	Way   OrderingWay
}

func (mf *ModelFilter) Subqueryload() Filter {
	panic(errors.NotImplementedf("Subqueryload"))
}

func (mf *ModelFilter) AllowSubqueryload() Filter {
	mf.subqueryload = true
	return mf
}

func (mf *ModelFilter) ShouldSubqueryload() bool {
	return mf.subqueryload
}

func (mf *ModelFilter) AddOption(opt RequestOption) Filter {
	panic(errors.NotImplementedf("AddOption"))
}

func (mf *ModelFilter) OrderBy(field string, way OrderingWay) Filter {
	panic(errors.NotImplementedf("OrderBy"))
}

func (mf *ModelFilter) SetOrderBy(field string, way OrderingWay) Filter {
	if mf.orderBy == nil {
		mf.orderBy = []*OrderBy{}
	}
	mf.orderBy = append(mf.orderBy, &OrderBy{field, way})
	return mf
}

func (mf *ModelFilter) GetOrderBy() []*OrderBy {
	return mf.orderBy[:]
}

func (mf *ModelFilter) AddOption_(opt RequestOption) {
	mf.options = append(mf.options, opt)
}

func (mf *ModelFilter) GetSelectOptions() []RequestOption {
	opts := []RequestOption{}
	for _, opt := range mf.options {
		switch opt {
		case RequestOptions.SelectForUpdate:
			opts = append(opts, opt)
		}
	}
	return opts
}
