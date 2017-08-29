package yaorm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
	"github.com/geoffreybauduin/yaorm/tools"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

type filterApplier struct {
	statement squirrel.SelectBuilder
	filter    yaormfilter.Filter
	tableName string
	dbp       DBProvider
}

type filterFieldApplier struct {
	statement   squirrel.SelectBuilder
	field       reflect.Value
	tableName   string
	tagData     []string
	dbFieldName string
	isJoining   bool
	leftJoin    bool
	dbp         DBProvider
}

func getTableNameFromFilter(f yaormfilter.Filter) string {
	table, err := GetTableByFilter(f)
	if err != nil {
		return ""
	}
	return table.Name()
}

func apply(statement squirrel.SelectBuilder, f yaormfilter.Filter, dbp DBProvider) squirrel.SelectBuilder {
	applier := &filterApplier{
		statement: statement,
		filter:    f,
		tableName: getTableNameFromFilter(f),
		dbp:       dbp,
	}
	applier.Apply()
	statement = applier.statement
	for _, option := range f.GetSelectOptions() {
		switch option {
		case yaormfilter.RequestOptions.SelectForUpdate:
			if dbp.CanSelectForUpdate() {
				statement = statement.Suffix(fmt.Sprintf(`FOR UPDATE OF %s`, dbp.EscapeValue(applier.tableName)))
			}
		}
	}
	return statement
}

func (a *filterApplier) Apply() {
	// We may have a pointer as parameter, we need to make sure we work with raw values
	underlyingFilter := tools.GetNonPtrValue(a.filter)
	st := underlyingFilter.Type()
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		dbFieldData, ok := field.Tag.Lookup("filter")
		if !ok || dbFieldData == "-" {
			// Skip
			continue
		}
		tagData := strings.Split(dbFieldData, ",")
		applier := &filterFieldApplier{
			field:     underlyingFilter.Field(i),
			tagData:   tagData,
			statement: a.statement,
			tableName: a.tableName,
			dbp:       a.dbp,
		}
		applier.Apply()
		a.statement = applier.statement
	}
}

func (a *filterFieldApplier) Apply() {
	a.setupFromTag()
	// If field is not nil, then apply the filter
	if a.fieldIsValid() {
		switch a.field.Kind() {
		case reflect.Ptr, reflect.Interface:
			a.applyPtr()
			break
		case reflect.Slice:
			a.applySlice()
			break
		default:
			panic(a.field.Kind())
		}
	}
}

func (a *filterFieldApplier) fieldIsValid() bool {
	return a.field.IsValid() && !a.field.IsNil()
}

func (a *filterFieldApplier) setupFromTag() {
	a.dbFieldName = a.tagData[0]
	a.isJoining = false
	a.leftJoin = false
	if len(a.tagData) == 4 && strings.Contains(a.tagData[1], "join") {
		a.isJoining = true
		if strings.Contains(a.tagData[1], "left") {
			a.leftJoin = true
		}
	}
}

func (a *filterFieldApplier) applyPtr() {
	valueFilter, ok := a.field.Interface().(yaormfilter.ValueFilter)
	if ok {
		a.statement = valueFilter.Apply(a.statement, a.dbp.EscapeValue(a.tableName), a.dbp.EscapeValue(a.dbFieldName))
	} else {
		structFilter, ok := a.field.Interface().(yaormfilter.Filter)
		if ok {
			tableName := getTableNameFromFilter(structFilter)
			if a.isJoining {
				tableName = fmt.Sprintf("%s_%s", a.tableName, tableName)
			} else {
				panic("what the fuck man, you're not joining")
			}
			a.applyFilter(structFilter, tableName)
		} else {
			// maybe it's another kind of filter ?
			panic(fmt.Errorf("I DONT KNOW HOW TO HANDLE THAT KIND OF FILTER %+v", a.field.Interface()))
		}
	}
}

func (a *filterFieldApplier) applySlice() {
	filterSlice := reflect.ValueOf(a.field.Interface())
	for idx := 0; idx < filterSlice.Len(); idx++ {
		f := filterSlice.Index(idx).Interface().(yaormfilter.Filter)
		tableName := getTableNameFromFilter(f)
		if a.isJoining {
			tableName = fmt.Sprintf("%s%d", tableName, idx)
		}
		a.applyFilter(f, tableName)
	}
}

func (a *filterFieldApplier) applyFilter(f yaormfilter.Filter, tableName string) {
	if !hasAnyFilter(f) {
		return
	}
	if a.isJoining {
		a.join(f, tableName)
	}
	filterApplier := &filterApplier{
		statement: a.statement,
		tableName: tableName,
		filter:    f,
		dbp:       a.dbp,
	}
	filterApplier.Apply()
	a.statement = filterApplier.statement
}

func (a *filterFieldApplier) join(f yaormfilter.Filter, tableAlias string) {
	joinCondition := fmt.Sprintf(
		`%s as %s on %s.%s = %s.%s`,
		a.dbp.EscapeValue(getTableNameFromFilter(f)),
		a.dbp.EscapeValue(tableAlias),
		a.dbp.EscapeValue(tableAlias),
		a.dbp.EscapeValue(a.tagData[2]),
		a.dbp.EscapeValue(a.tableName),
		a.dbp.EscapeValue(a.tagData[3]),
	)
	if !a.leftJoin {
		a.statement = a.statement.Join(joinCondition)
	} else {
		a.statement = a.statement.LeftJoin(joinCondition)
	}
}

func hasAnyFilter(f yaormfilter.Filter) bool {
	valueF := tools.GetNonPtrValue(f)
	if !valueF.IsValid() {
		return false
	}
	st := valueF.Type()
	for i := 0; i < st.NumField(); i++ {
		val, ok := st.Field(i).Tag.Lookup("filter")
		if !ok || val == "-" {
			continue
		}
		field := valueF.Field(i)
		if field.CanInterface() {
			_, ok := field.Interface().(yaormfilter.ValueFilter)
			if ok {
				return true
			}
			_, ok = field.Interface().(yaormfilter.Filter)
			if ok && hasAnyFilter(field.Interface().(yaormfilter.Filter)) {
				return true
			}
		} else {
			panic(field)
		}
	}
	return false
}
