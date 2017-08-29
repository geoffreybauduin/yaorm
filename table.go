package yaorm

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/loopfz/gadgeto/zesty/utils/rekordo"
	"github.com/geoffreybauduin/yaorm/tools"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/juju/errors"
)

var (
	// ErrTableEmpty is returned when a registered model has no exported fields
	ErrTableEmpty = errors.Errorf("Table is empty")
	// ErrTableNotFound is returned when a table cannot be found
	ErrTableNotFound = errors.NotFoundf("Table")
	// ErrDbNotFound is returned when a database cannot be found
	ErrDbNotFound = errors.NotFoundf("Database")
	tables        map[string]map[string]*Table
	tableByType   map[reflect.Type]*Table
	tableMutex    sync.RWMutex
)

func init() {
	tables = map[string]map[string]*Table{}
	tableByType = map[reflect.Type]*Table{}
}

// Table is the type hosting all the table characteristics
type Table struct {
	name                string
	dbname              string
	model               Model
	reflectedType       reflect.Type
	fields              []string
	tm                  *rekordo.TableModel
	filter              yaormfilter.Filter
	keys                []string
	fieldsByDbKey       map[string]int
	filterFieldsByDbKey map[string]int
}

// NewTable registers a new table
func NewTable(dbName, tableName string, model Model) *Table {
	table := &Table{
		name:                tableName,
		dbname:              dbName,
		model:               model,
		reflectedType:       reflect.TypeOf(model).Elem(),
		fieldsByDbKey:       map[string]int{},
		filterFieldsByDbKey: map[string]int{},
		keys:                []string{},
	}
	err := table.retrieveFields()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	table.tm = rekordo.RegisterTableModel(dbName, tableName, tools.GetNonPtrInterface(model))
	tableMutex.Lock()
	defer tableMutex.Unlock()
	if _, ok := tables[dbName]; !ok {
		tables[dbName] = map[string]*Table{}
	}
	tables[dbName][tableName] = table
	tableByType[reflect.TypeOf(model).Elem()] = table
	return table
}

// GetTable returns the table matching the parameters
func GetTable(dbName, tableName string) (*Table, error) {
	if _, ok := tables[dbName]; !ok {
		return nil, ErrDbNotFound
	}
	if _, ok := tables[dbName][tableName]; !ok {
		return nil, ErrTableNotFound
	}
	return tables[dbName][tableName], nil
}

// GetTableByModel returns the table registered for this model
func GetTableByModel(m Model) (*Table, error) {
	table, ok := tableByType[reflect.TypeOf(m).Elem()]
	if !ok {
		return nil, ErrTableNotFound
	}
	return table, nil
}

// GetTableByFilter returns the table using this filter
func GetTableByFilter(f yaormfilter.Filter) (*Table, error) {
	table, ok := tableByType[reflect.TypeOf(f).Elem()]
	if !ok {
		return nil, ErrTableNotFound
	}
	return table, nil
}

func (t *Table) retrieveFields() error {
	fields := []string{}
	st := reflect.TypeOf(reflect.ValueOf(t.model).Elem().Interface())
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		dbFieldData, ok := field.Tag.Lookup("db")
		if !ok || dbFieldData == "-" {
			continue
		}
		tagData := strings.Split(dbFieldData, ",")
		fields = append(fields, tagData[0])
		t.fieldsByDbKey[tagData[0]] = i
		if tagData[0] == "id" {
			t.keys = []string{"id"}
		}
	}
	t.fields = fields
	if len(fields) == 0 {
		return ErrTableEmpty
	}
	return nil
}

func (t *Table) retrieveFilterFields() {
	t.filterFieldsByDbKey = map[string]int{}
	st := reflect.TypeOf(reflect.ValueOf(t.filter).Elem().Interface())
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		dbFieldData, ok := field.Tag.Lookup("filter")
		if !ok || dbFieldData == "-" {
			continue
		}
		tagData := strings.Split(dbFieldData, ",")
		t.filterFieldsByDbKey[tagData[0]] = i
	}
}

func (t *Table) WithKeys(keys []string) *Table {
	t.keys = keys
	t.tm = t.tm.WithKeys(keys)
	return t
}

func (t *Table) WithFilter(f yaormfilter.Filter) *Table {
	if t.filter != nil {
		delete(tableByType, reflect.TypeOf(t.filter).Elem())
	}
	t.filter = f
	tableByType[reflect.TypeOf(t.filter).Elem()] = t
	t.retrieveFilterFields()
	return t
}

func (t *Table) WithAutoIncrement(v bool) *Table {
	t.tm = t.tm.WithAutoIncrement(v)
	return t
}

func (t *Table) WithSubqueryloading(fn SubqueryloadFunc, mapperField string) *Table {
	key := t.Name()
	if mapperField != "id" {
		key = fmt.Sprintf("%s_per_%s", key, mapperField)
	}
	subqueryloaders[key] = subqueryloader{
		fn:          fn,
		mapperField: mapperField,
	}
	return t
}

func (t Table) Fields() []string {
	return t.fields
}

func (t Table) FieldsForQuery(tableName string) []string {
	fields := []string{}
	for _, field := range t.fields {
		fields = append(fields, fmt.Sprintf(`"%s"."%s"`, tableName, field))
	}
	return fields
}

func (t Table) NameForQuery() string {
	return fmt.Sprintf(`"%s"`, t.Name())
}

func (t Table) Name() string {
	return t.name
}

func (t Table) Keys() []string {
	return t.keys
}

func (t Table) KeyFields() map[string]int {
	m := map[string]int{}
	for _, key := range t.Keys() {
		m[key] = t.fieldsByDbKey[key]
	}
	return m
}

func (t Table) FieldIndex(field string) int {
	idx, ok := t.fieldsByDbKey[field]
	if !ok {
		return -1
	}
	return idx
}

func (t Table) FilterFieldIndex(field string) int {
	idx, ok := t.filterFieldsByDbKey[field]
	if !ok {
		return -1
	}
	return idx
}

func (t Table) NewFilter() (yaormfilter.Filter, error) {
	f, ok := reflect.New(reflect.TypeOf(t.filter).Elem()).Interface().(yaormfilter.Filter)
	if !ok {
		return nil, errors.Errorf("Invalid filter format stored inside the Table instance")
	}
	return f, nil
}

func (t Table) NewModel() (Model, error) {
	m, ok := reflect.New(reflect.TypeOf(t.model).Elem()).Interface().(Model)
	if !ok {
		return nil, errors.Errorf("Invalid model format stored inside the Table instance")
	}
	return m, nil
}

func (t Table) NewSlice() (interface{}, error) {
	m, err := t.NewModel()
	if err != nil {
		return nil, err
	}
	sliceType := reflect.SliceOf(reflect.ValueOf(m).Type())
	slice := reflect.MakeSlice(sliceType, 1, 1)
	return slice.Interface(), nil
}

func (t Table) NewSlicePtr() (interface{}, error) {
	m, err := t.NewModel()
	if err != nil {
		return nil, err
	}
	sliceType := reflect.SliceOf(reflect.ValueOf(m).Type())
	slice := reflect.MakeSlice(sliceType, 1, 1)
	return reflect.New(slice.Type()).Interface(), nil
}
