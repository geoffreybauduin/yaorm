package testdata

import (
	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type PostType struct {
	yaorm.DatabaseModel
	PostID int64  `db:"post_id"`
	Type   string `db:"type"`
}

type PostTypeFilter struct {
	yaormfilter.ModelFilter
	FilterPostID yaormfilter.ValueFilter `filter:"post_id"`
}

func init() {
	yaorm.NewTable("test", "post_type", &PostType{}).WithFilter(&PostTypeFilter{}).WithKeys([]string{"post_id"}).WithAutoIncrement(false)
}

func (pt *PostType) Save() error {
	return yaorm.SaveWithPrimaryKeys(pt, []string{"post_id"})
}

func NewPostTypeFilter() *PostTypeFilter {
	return &PostTypeFilter{}
}

func (f *PostTypeFilter) PostID(v yaormfilter.ValueFilter) *PostTypeFilter {
	f.FilterPostID = v
	return f
}
