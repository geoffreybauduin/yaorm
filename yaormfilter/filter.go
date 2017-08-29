package yaormfilter

import "github.com/juju/errors"

// Filter is the interface stating the methods to implement to correctly filter on models
type Filter interface {
	Subqueryload() Filter
	ShouldSubqueryload() bool
	AddOption(opt RequestOption) Filter
	GetSelectOptions() []RequestOption
}

type RequestOption string

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
