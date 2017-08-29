package testdata

import (
	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type Post struct {
	yaorm.DatabaseModel
	ID         int64     `db:"id"`
	Subject    string    `db:"subject"`
	CategoryID int64     `db:"category_id"`
	Category   *Category `db:"-" filterload:"category,category_id"`
}

type PostFilter struct {
	yaormfilter.ModelFilter
	FilterID       yaormfilter.ValueFilter `filter:"id"`
	FilterSubject  yaormfilter.ValueFilter `filter:"subject"`
	FilterCategory yaormfilter.Filter      `filter:"category,join,id,category_id" filterload:"category"`
}

func init() {
	yaorm.NewTable("test", "post", &Post{}).WithFilter(&PostFilter{})
}

func NewPostFilter() *PostFilter {
	return &PostFilter{}
}

func (f *PostFilter) ID(v yaormfilter.ValueFilter) *PostFilter {
	f.FilterID = v
	return f
}

func (f *PostFilter) Subject(v yaormfilter.ValueFilter) *PostFilter {
	f.FilterSubject = v
	return f
}

func (f *PostFilter) Category(v yaormfilter.Filter) *PostFilter {
	if _, ok := v.(*CategoryFilter); !ok {
		panic("Not a CategoryFilter")
	}
	f.FilterCategory = v
	return f
}

func (f *PostFilter) Subqueryload() yaormfilter.Filter {
	f.AllowSubqueryload()
	return f
}
