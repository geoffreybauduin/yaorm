package yaorm

import (
	"github.com/geoffreybauduin/yaorm/tools"
	"github.com/go-gorp/gorp"
)

func isPGSQL(dbp DBProvider) bool {
	// This is shameful and sad, but I dont see how we can do this a better way
	v := tools.GetNonPtrValue(dbp.DB())
	// will work unless they change their stuff
	dbField := v.FieldByName("DB")
	realValue := tools.GetNonPtrValue(dbField.Interface())
	field := realValue.FieldByName("DbMap")
	s := field.Interface().(*gorp.DbMap)
	if _, ok := s.Dialect.(gorp.PostgresDialect); ok {
		return true
	}
	return false
}
