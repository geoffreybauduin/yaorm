package yaorm

import (
	"fmt"
	"reflect"
	"github.com/geoffreybauduin/yaorm/tools"
)

type SubqueryloadFunc func(dbp DBProvider, ids []interface{}) (interface{}, error)

type subqueryloader struct {
	// fn to retrieve the data
	fn SubqueryloadFunc
	// field on loaded model from which to retrieve the associated fk from original model
	mapperField string
}

var (
	subqueryloaders = map[string]subqueryloader{}
)

func subqueryload(dbp DBProvider, m map[string]map[interface{}][]reflect.Value) (map[string][]Model, error) {
	models := map[string][]Model{}
	for model, data := range m {
		modelData, err := subqueryloadForModel(dbp, model, data)
		if err != nil {
			return models, err
		}
		models[model] = modelData
	}
	return models, nil
}

// avoid hammering the database and selecting more than x rows at a time ?
// pg max is 65536
const loadStep = 1000

func subqueryloadForModel(dbp DBProvider, model string, data map[interface{}][]reflect.Value) ([]Model, error) {
	d := []Model{}
	if _, ok := subqueryloaders[model]; !ok {
		return d, fmt.Errorf("Subqueryload model %s not defined yet", model)
	}
	ids := []interface{}{}
	for id, _ := range data {
		ids = append(ids, id)
	}
	for factor := 0; factor*loadStep < len(ids); factor++ {
		var subset []interface{}
		currentStep := factor * loadStep
		nextStep := (factor + 1) * loadStep
		if nextStep < len(ids) {
			subset = ids[currentStep:nextStep]
		} else {
			subset = ids[currentStep:]
		}
		models, err := subqueryloaders[model].fn(dbp, subset)
		if err != nil {
			return d, err
		}
		modelsSlice := reflect.ValueOf(models)
		for i := 0; i < modelsSlice.Len(); i++ {
			m := modelsSlice.Index(i)
			table, err := GetTableByModel(m.Interface().(Model))
			if err != nil {
				return nil, err
			}
			field := tools.GetNonPtrValue(m.Interface()).Field(table.FieldIndex(subqueryloaders[model].mapperField))
			fk := reflect.Indirect(field).Interface()
			for _, v := range data[fk] {
				setOnReceiver(v, m.Interface())
			}
			d = append(d, m.Interface().(Model))
		}
	}
	return d, nil
}

func setOnReceiver(v reflect.Value, value interface{}) {
	ind := reflect.Indirect(v)
	valueToAdd := reflect.ValueOf(value)
	switch ind.Kind() {
	case reflect.Slice:
		newVal := reflect.Append(ind, valueToAdd)
		ind.Set(newVal)
		break
	case reflect.Ptr:
		ind.Set(valueToAdd)
		break
	default:
		panic(fmt.Errorf("Wow, why do you have %+v receiver ?", ind.Kind()))
	}
}
