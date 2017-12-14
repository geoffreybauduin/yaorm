package testdata

import (
	"time"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type InvalidPostTag struct {
	yaorm.DatabaseModel
	PostID    int64     `db:"post_id"`
	TagID     int64     `db:"tag_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type InvalidPostTagFilter struct {
	yaormfilter.ModelFilter
	FilterPostID yaormfilter.ValueFilter `filter:"post_id"`
}

func init() {
	yaorm.NewTable("test", "invalid_post_tag", &InvalidPostTag{}).WithFilter(NewInvalidPostTagFilter()).WithKeys([]string{"post_id", "tag_id"}).WithAutoIncrement(false)
}

func (pt *InvalidPostTag) DBHookBeforeInsert() error {
	now := time.Now()
	pt.CreatedAt = now
	pt.UpdatedAt = now
	return nil
}

func (pt *InvalidPostTag) DBHookBeforeUpdate() error {
	now := time.Now()
	pt.UpdatedAt = now
	return nil
}

func (pt *InvalidPostTag) Save() error {
	return yaorm.GenericSave(pt)
}

func NewInvalidPostTagFilter() *InvalidPostTagFilter {
	return &InvalidPostTagFilter{}
}

func (f *InvalidPostTagFilter) PostID(v yaormfilter.ValueFilter) *InvalidPostTagFilter {
	f.FilterPostID = v
	return f
}
