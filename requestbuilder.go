package yaorm

import (
	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
)

var (
	generator = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func buildSelect(m Model) (squirrel.SelectBuilder, error) {
	table, err := GetTableByModel(m)
	if err != nil {
		return squirrel.SelectBuilder{}, err
	}
	fields := table.FieldsForQuery(table.Name())

	return generator.Select(fields...).From(table.NameForQuery()), nil
}
