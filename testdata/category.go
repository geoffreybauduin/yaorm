package testdata

import (
	"time"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type Category struct {
	yaorm.DatabaseModel
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CategoryFilter struct {
	yaormfilter.ModelFilter
	FilterID   yaormfilter.ValueFilter `filter:"id"`
	FilterName yaormfilter.ValueFilter `filter:"name"`
}

func init() {
	yaorm.NewTable("test", "category", &Category{}).WithFilter(&CategoryFilter{}).WithSubqueryloading(
		func(dbp yaorm.DBProvider, ids []interface{}) (interface{}, error) {
			return yaorm.GenericSelectAll(dbp, NewCategoryFilter().ID(yaormfilter.In(ids...)))
		},
		"id",
	)
}

func (c *Category) DBHookBeforeInsert() error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Category) DBHookBeforeUpdate() error {
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Category) Load(dbp yaorm.DBProvider) error {
	return yaorm.GenericSelectOneFromModel(dbp, c)
}

func (c *Category) Save() error {
	return yaorm.GenericSave(c)
}

func NewCategoryFilter() *CategoryFilter {
	return &CategoryFilter{}
}

func (f *CategoryFilter) ID(v yaormfilter.ValueFilter) *CategoryFilter {
	f.FilterID = v
	return f
}

func (f *CategoryFilter) Name(v yaormfilter.ValueFilter) *CategoryFilter {
	f.FilterName = v
	return f
}

func (f *CategoryFilter) Subqueryload() yaormfilter.Filter {
	f.AllowSubqueryload()
	return f
}
