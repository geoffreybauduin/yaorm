package yaorm

import (
	"fmt"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
	"github.com/geoffreybauduin/yaorm/tools"
)

func buildSelect(dbp DBProvider, m Model) (squirrel.SelectBuilder, error) {
	table, err := GetTableByModel(m)
	if err != nil {
		return squirrel.SelectBuilder{}, err
	}
	fields := table.Fields()
	var f []string
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

func buildInsert(dbp DBProvider, m Model) (squirrel.InsertBuilder, error) {
	table, err := GetTableByModel(m)
	if err != nil {
		return squirrel.InsertBuilder{}, err
	}
	reflectedM := tools.GetNonPtrValue(m)
	stmt := dbp.getStatementGenerator().Insert(table.Name()).Columns(table.FieldsWithoutPK()...)
	var values []interface{}
	for _, field := range table.FieldsWithoutPK() {
		values = append(values, reflectedM.Field(table.FieldIndex(field)).Interface())
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
		stmt = stmt.Set(field, reflectedM.Field(table.FieldIndex(field)).Interface())
	}
	for pk, idx := range table.KeyFields() {
		stmt = stmt.Where(squirrel.Eq{pk: reflectedM.Field(idx).Interface()})
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
		stmt = stmt.Where(squirrel.Eq{pk: reflectedM.Field(idx).Interface()})
	}
	return stmt, nil
}
