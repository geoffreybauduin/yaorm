package testdata

import (
	"time"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type PostTag struct {
	yaorm.DatabaseModel
	PostID    int64     `db:"post_id"`
	TagID     int64     `db:"tag_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PostTagFilter struct {
	yaormfilter.ModelFilter
	FilterPostID yaormfilter.ValueFilter `filter:"post_id"`
	FilterTagID  yaormfilter.ValueFilter `filter:"tag_id"`
}

func init() {
	yaorm.NewTable("test", "post_tag", &PostTag{}).WithFilter(NewPostTagFilter()).WithKeys([]string{"post_id", "tag_id"}).WithAutoIncrement(false)
}

func (pt *PostTag) DBHookBeforeInsert() error {
	now := time.Now()
	pt.CreatedAt = now
	pt.UpdatedAt = now
	return nil
}

func (pt *PostTag) DBHookBeforeUpdate() error {
	now := time.Now()
	pt.UpdatedAt = now
	return nil
}

func NewPostTagFilter() *PostTagFilter {
	return &PostTagFilter{}
}

func (f *PostTagFilter) PostID(v yaormfilter.ValueFilter) *PostTagFilter {
	f.FilterPostID = v
	return f
}

func (f *PostTagFilter) TagID(v yaormfilter.ValueFilter) *PostTagFilter {
	f.FilterTagID = v
	return f
}
