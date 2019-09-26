package yaorm

import (
	"fmt"
	"reflect"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
	"github.com/geoffreybauduin/yaorm/tools"
)

type buildSelectColumns struct {
	loadColumns     []string
	dontLoadColumns []string
}

func (c *buildSelectColumns) reduce(columns []string) []string {
	emptyLoad := len(c.loadColumns) == 0
	emptyDontLoad := len(c.dontLoadColumns) == 0
	if emptyLoad && emptyDontLoad {
		return columns
	}

	// One cache for both load & dontLoad
	cache := make(map[string]bool, len(c.loadColumns)+len(c.dontLoadColumns))
	for _, col := range c.loadColumns {
		cache[col] = true
	}
	for _, col := range c.dontLoadColumns {
		cache[col] = false
	}

	// Do not reuse columns slice to avoid polluting its backend array
	forceExclude := !emptyDontLoad && emptyLoad
	var ok, exists bool
	res := make([]string, 0, len(columns))
	for _, f := range columns {
		ok, exists = cache[f]
		if ok || (!exists && forceExclude) {
			res = append(res, f)
		}
	}
	return res
}

func buildSelect(dbp DBProvider, m Model, sf buildSelectColumns) (squirrel.SelectBuilder, error) {
	table, err := GetTableByModel(m)
	if err != nil {
		return squirrel.SelectBuilder{}, err
	}
	fields := sf.reduce(table.Fields())
	f := make([]string, 0, len(fields))
	for _, field := range fields {
		f = append(f, fmt.Sprintf(`%s.%s`, dbp.EscapeValue(table.Name()), dbp.EscapeValue(field)))
	}

	return dbp.getStatementGenerator().Select(f...).From(
		fmt.Sprintf("%s AS %s", table.NameForQuery(dbp), dbp.EscapeValue(table.Name())),
	), nil
}

func buildCount(dbp DBProvider, table *Table) (squirrel.SelectBuilder, error) {
	return dbp.getStatementGenerator().Select("COUNT(*) AS count").From(dbp.EscapeValue(table.Name())), nil
}

func getNiceArgumentFormatted(v interface{}) string {
	if v == nil {
		return "nil"
	}
	value := reflect.ValueOf(v)
	if tools.IsZeroValue(value) {
		return fmt.Sprintf("%v", v)
	}
	switch value.Kind() {
	case reflect.Ptr:
		return getNiceArgumentFormatted(reflect.Indirect(value).Interface())

	}
	return fmt.Sprintf("%v", v)
}

func buildInsert(dbp DBProvider, m Model) (squirrel.InsertBuilder, error) {
	table, err := GetTableByModel(m)
	if err != nil {
		return squirrel.InsertBuilder{}, err
	}
	reflectedM := tools.GetNonPtrValue(m)
	stmt := dbp.getStatementGenerator().Insert(table.Name()).Columns(table.FieldsWithoutPK()...)
	var values []interface{}
	for _, field := range table.FieldsWithoutPK() {
		values = append(values, getNiceArgumentFormatted(reflectedM.Field(table.FieldIndex(field)).Interface()))
	}
	stmt = stmt.Values(values...)
	return stmt, nil
}

func buildUpdate(dbp DBProvider, m Model) (squirrel.UpdateBuilder, error) {
	table, err := GetTableByModel(m)
	if err != nil {
		return squirrel.UpdateBuilder{}, err
	}
	reflectedM := tools.GetNonPtrValue(m)
	stmt := dbp.getStatementGenerator().Update(table.Name())
	for _, field := range table.FieldsWithoutPK() {
		stmt = stmt.Set(field, getNiceArgumentFormatted(reflectedM.Field(table.FieldIndex(field)).Interface()))
	}
	for pk, idx := range table.KeyFields() {
		stmt = stmt.Where(squirrel.Eq{pk: tools.GetNonPtrInterface(reflectedM.Field(idx).Interface())})
	}
	return stmt, nil
}

func buildDelete(dbp DBProvider, m Model) (squirrel.DeleteBuilder, error) {
	table, err := GetTableByModel(m)
	if err != nil {
		return squirrel.DeleteBuilder{}, err
	}
	reflectedM := tools.GetNonPtrValue(m)
	stmt := dbp.getStatementGenerator().Delete(table.Name())
	for pk, idx := range table.KeyFields() {
		stmt = stmt.Where(squirrel.Eq{pk: getNiceArgumentFormatted(reflectedM.Field(idx).Interface())})
	}
	return stmt, nil
}
