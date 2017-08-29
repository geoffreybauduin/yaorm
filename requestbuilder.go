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
	f := []string{}
	for _, field := range fields {
		f = append(f, fmt.Sprintf(`%s.%s`, dbp.EscapeValue(table.Name()), dbp.EscapeValue(field)))
	}

	return dbp.getStatementGenerator().Select(f...).From(dbp.EscapeValue(table.Name())), nil
}
