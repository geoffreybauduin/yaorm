package testdata

import (
	"fmt"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type Post struct {
	yaorm.DatabaseModel
	ID           int64           `db:"id"`
	Subject      string          `db:"subject"`
	CategoryID   int64           `db:"category_id"`
	Category     *Category       `db:"-" filterload:"category,category_id"`
	ParentPostID int64           `db:"parent_post_id"`
	ChildrenPost []*Post         `db:"-" filterload:"post,id,parent_post_id"`
	Metadata     []*PostMetadata `db:"-" filterload:"post_metadata,id,post_id"`
}

type PostFilter struct {
	yaormfilter.ModelFilter
	FilterID           yaormfilter.ValueFilter `filter:"id"`
	FilterParentPostID yaormfilter.ValueFilter `filter:"parent_post_id"`
	FilterSubject      yaormfilter.ValueFilter `filter:"subject"`
	FilterCategory     yaormfilter.Filter      `filter:"category,join,id,category_id" filterload:"category"`
	FilterChildren     yaormfilter.Filter      `filter:"post,join,parent_post_id,id" filterload:"post"`
	FilterMetadata     []yaormfilter.Filter    `filter:"post_metadata,join,post_id,id" filterload:"post_metadata"`
}

func init() {
	yaorm.NewTable("test", "post", &Post{}).WithFilter(&PostFilter{}).WithSubqueryloading(
		func(dbp yaorm.DBProvider, ids []interface{}) (interface{}, error) {
			return yaorm.GenericSelectAll(dbp, NewPostFilter().ParentPostID(yaormfilter.In(ids...)))
		}, "parent_post_id",
	)
}

func NewPostFilter() *PostFilter {
	return &PostFilter{}
}

func (f *PostFilter) ID(v yaormfilter.ValueFilter) *PostFilter {
	f.FilterID = v
	return f
}

func (f *PostFilter) ParentPostID(v yaormfilter.ValueFilter) *PostFilter {
	f.FilterParentPostID = v
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

func (f *PostFilter) ChildrenPosts(v yaormfilter.Filter) *PostFilter {
	if _, ok := v.(*PostFilter); !ok {
		panic("Not a PostFilter")
	}
	f.FilterChildren = v
	return f
}

func (f *PostFilter) Subqueryload() yaormfilter.Filter {
	f.AllowSubqueryload()
	return f
}

func (f *PostFilter) Metadata(metadata ...yaormfilter.Filter) *PostFilter {
	if f.FilterMetadata == nil {
		f.FilterMetadata = make([]yaormfilter.Filter, 0)
	}
	for _, md := range metadata {
		if _, ok := md.(*PostMetadataFilter); !ok {
			panic(fmt.Errorf("filter %v is not a PostMetadataFilter", md))
		}
		f.FilterMetadata = append(f.FilterMetadata, md)
	}
	return f
}
