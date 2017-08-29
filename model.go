package yaorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/juju/errors"

	"github.com/geoffreybauduin/yaorm/tools"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

// Model is the interface every model should implement from
type Model interface {
	GetDBP() DBProvider
	SetDBP(dbp DBProvider)
	Save() error
	Load(dbp DBProvider) error
	Delete() error
	DBHookBeforeInsert() error
	DBHookBeforeUpdate() error
}

// DatabaseModel is the struct every model should compose
type DatabaseModel struct {
	dbp DBProvider `db:"-"`
}

func (dm *DatabaseModel) SetDBP(dbp DBProvider) {
	dm.dbp = dbp
}

func (dm *DatabaseModel) GetDBP() DBProvider {
	return dm.dbp
}

func (dm *DatabaseModel) Save() error {
	return errors.NotImplementedf("Save")
}

func (dm *DatabaseModel) Load(dbp DBProvider) error {
	return errors.NotImplementedf("Load")
}

func (dm *DatabaseModel) Delete() error {
	return errors.NotImplementedf("Delete")
}

func (dm *DatabaseModel) DBHookBeforeInsert() error {
	return nil
}

func (dm *DatabaseModel) DBHookBeforeUpdate() error {
	return nil
}

// GenericSelectOne selects one row in the database
// panics if filter or dbp is nil
func GenericSelectOne(dbp DBProvider, filter yaormfilter.Filter) (Model, error) {
	table, err := GetTableByFilter(filter)
	if err != nil {
		return nil, err
	}
	m, err := table.NewModel()
	if err != nil {
		return nil, err
	}
	err = GenericSelectOneWithModel(dbp, filter, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GenericSelectOneFromModel selects one row in the database from a model value
// panics if filter or dbp is nil
func GenericSelectOneFromModel(dbp DBProvider, m Model) error {
	table, err := GetTableByModel(m)
	if err != nil {
		return err
	}
	filter, err := table.NewFilter()
	if err != nil {
		return err
	}
	valueM := tools.GetNonPtrValue(m)
	typeM := valueM.Type()
	for i := 0; i < typeM.NumField(); i++ {
		tagValue, ok := typeM.Field(i).Tag.Lookup("db")
		if !ok || tagValue == "-" {
			continue
		}
		tagData := strings.Split(tagValue, ",")
		value := valueM.Field(i)
		if !tools.IsZeroValue(value) {
			idx := table.FilterFieldIndex(tagData[0])
			if idx < 0 {
				return errors.Errorf("Cannot find field with filter tag '%s' in filter %T", tagData[0], filter)
			}
			filterToApply := reflect.ValueOf(yaormfilter.Equals(value.Interface()))
			tools.GetNonPtrValue(filter).Field(idx).Set(filterToApply)
		}
	}
	return GenericSelectOneWithModel(dbp, filter, m)
}

// GenericSelectOneWithModel selects one row in the database providing the destination model directly
// panics if filter or dbp is nil
func GenericSelectOneWithModel(dbp DBProvider, filter yaormfilter.Filter, m Model) error {
	statement, err := buildSelect(m)
	if err != nil {
		return err
	}
	statement = apply(statement, filter, dbp)
	query, params, err := statement.ToSql()
	if err != nil {
		return err
	}
	err = dbp.DB().SelectOne(m, query, params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NotFoundf(formatTablenameFromError(m))
		}
		return err
	}
	m.SetDBP(dbp)
	if filter != nil {
		err = finishSelect(dbp, m, filter)
	}
	return err
}

// GenericSelectAll selects all rows in the database
// panics if filter or dbp is nil
func GenericSelectAll(dbp DBProvider, filter yaormfilter.Filter) ([]Model, error) {
	table, err := GetTableByFilter(filter)
	if err != nil {
		return nil, err
	}
	m, err := table.NewModel()
	if err != nil {
		return nil, err
	}
	statement, err := buildSelect(m)
	if err != nil {
		return nil, err
	}
	sm, _ := table.NewSlicePtr()
	statement = apply(statement, filter, dbp)
	query, params, err := statement.ToSql()
	if err != nil {
		return nil, err
	}
	_, err = dbp.DB().Select(sm, query, params...)
	if err != nil {
		panic(err)
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundf(formatTablenameFromError(m))
		}
		return nil, err
	}
	models := []Model{}
	smValue := tools.GetNonPtrValue(sm)
	for i := 0; i < smValue.Len(); i++ {
		m := smValue.Index(i).Interface().(Model)
		m.SetDBP(dbp)
		models = append(models, m)
	}
	if filter != nil {
		err = finishSelect(dbp, models, filter)
	}
	return models, err
}

// GenericSave updates or inserts the provided model in the database
// panics if model is nil or not linked to dbp
func GenericSave(m Model) error {
	table, err := GetTableByModel(m)
	if err != nil {
		return err
	}
	// loop over primary keys
	// if they are both to zero value, then it means it's an insert
	// if they are both set to values, then it's an update
	// otherwise, scream
	keys := table.KeyFields()
	if len(keys) == 1 {
		return genericSaveOnePK(m, keys)
	}
	return genericSaveMultiplePK(table, m, keys)
}

func genericSaveMultiplePK(table *Table, m Model, keys map[string]int) error {
	f, err := table.NewFilter()
	if err != nil {
		return err
	}
	reflectedModel := tools.GetNonPtrValue(m)
	for fieldName, idx := range keys {
		field := reflectedModel.Field(idx)
		if tools.IsZeroValue(field) {
			return errors.Errorf("Cannot save this model: one Primary Key has zero value")
		}
		filterIdx := table.FilterFieldIndex(fieldName)
		tools.GetNonPtrValue(f).Field(filterIdx).Set(reflect.ValueOf(yaormfilter.Equals(field.Interface())))
	}
	_, err = GenericSelectOne(m.GetDBP(), f)
	if err != nil && errors.IsNotFound(err) {
		return GenericInsert(m)
	} else if err == nil {
		return GenericUpdate(m)
	}
	return err
}

func genericSaveOnePK(m Model, keys map[string]int) error {
	isZero := true
	reflectedModel := tools.GetNonPtrValue(m)
	for _, idx := range keys {
		field := reflectedModel.Field(idx)
		if !tools.IsZeroValue(field) {
			isZero = false
			break
		}
	}
	if isZero {
		return GenericInsert(m)
	}
	return GenericUpdate(m)
}

// GenericUpdate updates the provided model in the database
// panics if model is nil or not linked to dbp
func GenericUpdate(m Model) error {
	err := m.DBHookBeforeUpdate()
	if err != nil {
		return err
	}
	_, err = m.GetDBP().DB().Update(m)
	return err
}

// GenericInsert inserts the provided model in the database
// panics if model is nil or not linked to dbp
func GenericInsert(m Model) error {
	err := m.DBHookBeforeInsert()
	if err != nil {
		return err
	}
	err = m.GetDBP().DB().Insert(m)
	return err
}

func formatTablenameFromError(m Model) string {
	table, err := GetTableByModel(m)
	if err != nil {
		return ""
	}
	tableName := table.Name()
	tableName = strings.Replace(tableName, "_", " ", -1)
	words := strings.Split(tableName, " ")
	for i := 0; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

func finishSelect(dbp DBProvider, m interface{}, f yaormfilter.Filter) error {
	fkPerModel := map[string]map[interface{}][]reflect.Value{}
	valueF := reflect.Indirect(reflect.ValueOf(f))
	if !valueF.IsValid() {
		return nil
	}
	st := valueF.Type()
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		dbFieldData, ok := field.Tag.Lookup("filterload")
		if !ok || dbFieldData == "-" {
			// Skip
			continue
		}
		filterFound := valueF.Field(i)
		if !filterFound.IsNil() && filterFound.Interface().(yaormfilter.Filter).ShouldSubqueryload() {
			fkPerModel = feedFkPerModel(m, fkPerModel, dbFieldData)
		}
	}
	d, err := subqueryload(dbp, fkPerModel)
	if err != nil {
		return err
	}
	for model, data := range d {
		realModel := strings.Split(model, "_per_")[0]
		for i := 0; i < st.NumField(); i++ {
			field := st.Field(i)
			dbFieldData, ok := field.Tag.Lookup("filterload")
			if !ok || dbFieldData != realModel {
				// Skip
				continue
			}
			filterFound := valueF.Field(i)
			if !filterFound.IsNil() && filterFound.Interface().(yaormfilter.Filter).ShouldSubqueryload() {
				err = finishSelect(dbp, data, filterFound.Interface().(yaormfilter.Filter))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func feedFkPerModel(m interface{}, fkPerModel map[string]map[interface{}][]reflect.Value, dbFieldData string) map[string]map[interface{}][]reflect.Value {
	s := reflect.Indirect(reflect.ValueOf(m))
	if s.Kind() == reflect.Slice {
		for i := 0; i < s.Len(); i++ {
			fkPerModel = feedFkPerModel(s.Index(i).Interface(), fkPerModel, dbFieldData)
		}
	} else {
		idx, tag := getFieldInModel(m.(Model), "filterload", dbFieldData)
		if idx > -1 {
			tagData := strings.Split(tag, ",")
			if len(tagData) == 3 {
				// it's a reverse filter
				tagData[0] = fmt.Sprintf("%s_per_%s", tagData[0], tagData[2])
			}
			idxFk, _ := getFieldInModel(m.(Model), "db", tagData[1])
			if idxFk > -1 {
				fkReflectedvalue := s.Field(idxFk)
				if !tools.IsZeroValue(fkReflectedvalue) {
					fkValue := reflect.Indirect(fkReflectedvalue).Interface()
					if _, ok := fkPerModel[tagData[0]]; !ok {
						fkPerModel[tagData[0]] = map[interface{}][]reflect.Value{}
					}
					if _, ok := fkPerModel[tagData[0]][fkValue]; !ok {
						fkPerModel[tagData[0]][fkValue] = []reflect.Value{}
					}
					fkPerModel[tagData[0]][fkValue] = append(fkPerModel[tagData[0]][fkValue], s.Field(idx).Addr())
				}
			} else {
				panic("Cannot find fk in model struct")
			}
		} else {
			panic(fmt.Errorf("Cannot find tag filterload:%v in model struct", dbFieldData))
		}
	}
	return fkPerModel
}

func getFieldInModel(m Model, tag, equals string) (int, string) {
	st := reflect.Indirect(reflect.ValueOf(m)).Type()
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		dbFieldData, ok := field.Tag.Lookup(tag)
		if !ok {
			continue
		}
		tagData := strings.Split(dbFieldData, ",")
		if tagData[0] == equals {
			return i, dbFieldData
		}
	}
	return -1, ""
}
