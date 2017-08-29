package testdata

import (
	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type Tag struct {
	yaorm.DatabaseModel
	ID  int64  `db:"id"`
	Tag string `db:"tag"`
}

type TagFilter struct {
	yaormfilter.ModelFilter
	FilterID  yaormfilter.ValueFilter `filter:"id"`
	FilterTag yaormfilter.ValueFilter `filter:"tag"`
}

func init() {
	yaorm.NewTable("test", "tag", &Tag{}).WithFilter(NewTagFilter())
}

func NewTagFilter() *TagFilter {
	return &TagFilter{}
}

func (f *TagFilter) ID(v yaormfilter.ValueFilter) *TagFilter {
	f.FilterID = v
	return f
}

func (f *TagFilter) Tag(v yaormfilter.ValueFilter) *TagFilter {
	f.FilterTag = v
	return f
}
