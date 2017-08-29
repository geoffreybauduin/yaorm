package testdata

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/geoffreybauduin/yaorm"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func SetupTestDatabase(name string) (func(), error) {
	tmpFile := fmt.Sprintf("/tmp/yaorm_%s_%d.sqlite", name, rand.Int())
	err := yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             name,
		DSN:              tmpFile,
		System:           yaorm.DatabaseSqlite3,
		AutoCreateTables: true,
	})
	if err != nil {
		return nil, err
	}
	return func() {
		err := yaorm.UnregisterDB(name)
		if err != nil {
			panic(err)
		}
		os.Remove(tmpFile)
	}, nil
}
