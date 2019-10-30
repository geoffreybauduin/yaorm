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
	Limit(limit uint64) Filter
	Offset(offset uint64) Filter
	GetLimit() (bool, uint64)
	GetOffset() (bool, uint64)
	LoadColumns(columns ...string)
	GetLoadColumns() []string
	DontLoadColumns(columns ...string)
	GetDontLoadColumns() []string
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
	SelectDistinct  RequestOption
	LeftJoin        RequestOption
}{
	SelectForUpdate: "SelectForUpdate",
	SelectDistinct:  "SelectDistinct",
	LeftJoin:        "LeftJoin",
}

// ModelFilter is the struct every filter should compose
type ModelFilter struct {
	subqueryload    bool
	options         []RequestOption
	orderBy         []*OrderBy
	shouldLimit     bool
	limit           uint64
	shouldOffset    bool
	offset          uint64
	loadColumns     []string
	dontLoadColumns []string
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
		case RequestOptions.SelectForUpdate, RequestOptions.SelectDistinct:
			opts = append(opts, opt)
		}
	}
	return opts
}

func (mf *ModelFilter) Limit(limit uint64) Filter {
	panic(errors.NotImplementedf("Limit"))
}

func (mf *ModelFilter) Offset(limit uint64) Filter {
	panic(errors.NotImplementedf("Offset"))
}

func (mf *ModelFilter) SetLimit(limit uint64) {
	mf.shouldLimit = true
	mf.limit = limit
}

func (mf *ModelFilter) SetOffset(offset uint64) {
	mf.shouldOffset = true
	mf.offset = offset
}

func (mf *ModelFilter) GetLimit() (bool, uint64) {
	return mf.shouldLimit, mf.limit
}

func (mf *ModelFilter) GetOffset() (bool, uint64) {
	return mf.shouldOffset, mf.offset
}

func (mf *ModelFilter) LoadColumns(columns ...string) {
	mf.loadColumns = append(mf.loadColumns, columns...)
}

func (mf *ModelFilter) GetLoadColumns() []string {
	return mf.loadColumns[:len(mf.loadColumns):len(mf.loadColumns)]
}

func (mf *ModelFilter) DontLoadColumns(columns ...string) {
	mf.dontLoadColumns = append(mf.dontLoadColumns, columns...)
}

func (mf *ModelFilter) GetDontLoadColumns() []string {
	return mf.dontLoadColumns[:len(mf.dontLoadColumns):len(mf.dontLoadColumns)]
}
