package yaorm

import "github.com/go-gorp/gorp"

//TypeConverter defines conversion from db to golang
type TypeConverter struct{}

//ToDb converts to database
func (tc TypeConverter) ToDb(val interface{}) (interface{}, error) {
	return val, nil
}

//FromDb converts to golang
func (tc TypeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	return gorp.CustomScanner{}, false
}
