package testdata

import (
	"time"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type TwoI struct {
	yaorm.DatabaseModel
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type TwoIFilter struct {
	yaormfilter.ModelFilter
	FilterID   yaormfilter.ValueFilter `filter:"id"`
	FilterName yaormfilter.ValueFilter `filter:"name"`
}

func init() {
	yaorm.NewTable("test", "2i", &TwoI{}).WithFilter(&TwoIFilter{}).WithSubqueryloading(
		func(dbp yaorm.DBProvider, ids []interface{}) (interface{}, error) {
			return yaorm.GenericSelectAll(dbp, NewTwoIFilter().ID(yaormfilter.In(ids...)))
		},
		"id",
	)
}

func (c *TwoI) DBHookBeforeInsert() error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *TwoI) DBHookBeforeUpdate() error {
	c.UpdatedAt = time.Now()
	return nil
}

func (c *TwoI) Load(dbp yaorm.DBProvider) error {
	return yaorm.GenericSelectOneFromModel(dbp, c)
}

func (c *TwoI) Save() error {
	return yaorm.GenericSave(c)
}

func NewTwoIFilter() *TwoIFilter {
	return &TwoIFilter{}
}

func (f *TwoIFilter) ID(v yaormfilter.ValueFilter) *TwoIFilter {
	f.FilterID = v
	return f
}

func (f *TwoIFilter) Name(v yaormfilter.ValueFilter) *TwoIFilter {
	f.FilterName = v
	return f
}

func (f *TwoIFilter) Subqueryload() yaormfilter.Filter {
	f.AllowSubqueryload()
	return f
}

func (f *TwoIFilter) OrderBy(field string, way yaormfilter.OrderingWay) yaormfilter.Filter {
	f.SetOrderBy(field, way)
	return f
}
