package yaorm

import (
	"fmt"
	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
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
		fmt.Sprintf("%s AS %s", table.NameForQuery(dbp), table.Name()),
	), nil
}

func buildCount(dbp DBProvider, table *Table) (squirrel.SelectBuilder, error) {
	return dbp.getStatementGenerator().Select("COUNT(*) AS count").From(dbp.EscapeValue(table.Name())), nil
}
