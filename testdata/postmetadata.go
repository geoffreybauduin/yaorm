package testdata

import (
	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type PostMetadata struct {
	yaorm.DatabaseModel
	PostID int64  `db:"post_id"`
	Key    string `db:"key"`
	Value  string `db:"value"`
}

type PostMetadataFilter struct {
	yaormfilter.ModelFilter
	FilterPostID yaormfilter.ValueFilter `filter:"post_id"`
	FilterKey    yaormfilter.ValueFilter `filter:"key"`
}

func (pm *PostMetadata) Save() error {
	return yaorm.GenericSave(pm)
}

func init() {
	yaorm.NewTable("test", "post_metadata", &PostMetadata{}).WithFilter(NewPostMetadataFilter()).WithKeys([]string{"post_id", "key"}).WithAutoIncrement(false).WithSubqueryloading(
		func(dbp yaorm.DBProvider, ids []interface{}) (interface{}, error) {
			return yaorm.GenericSelectAll(dbp, NewPostMetadataFilter().PostID(yaormfilter.In(ids...)))
		}, "post_id",
	)
}

func NewPostMetadataFilter() *PostMetadataFilter {
	return &PostMetadataFilter{}
}

func (f *PostMetadataFilter) PostID(v yaormfilter.ValueFilter) *PostMetadataFilter {
	f.FilterPostID = v
	return f
}

func (f *PostMetadataFilter) Key(v yaormfilter.ValueFilter) *PostMetadataFilter {
	f.FilterKey = v
	return f
}

func (f *PostMetadataFilter) Subqueryload() yaormfilter.Filter {
	f.AllowSubqueryload()
	return f
}
