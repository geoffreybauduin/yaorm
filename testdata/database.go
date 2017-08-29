package testdata

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/geoffreybauduin/yaorm"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var (
	tables = []string{"category", "post", "post_tag", "tag"}
)

func SetupTestDatabase(name string) (func(), error) {
	switch os.Getenv("DB") {
	case "postgres":
		return setupPostgres(name)
	case "mysql":
		return setupMysql(name)
	default:
		return setupSqlite(name)
	}
}

func setupSqlite(name string) (func(), error) {
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

func setupPostgres(name string) (func(), error) {
	err := yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             name,
		DSN:              os.Getenv("DSN"),
		System:           yaorm.DatabasePostgreSQL,
		AutoCreateTables: true,
	})
	if err != nil {
		return nil, err
	}
	return func() {
		dbp, err := yaorm.NewDBProvider(name)
		if err != nil {
			panic(err)
		}
		for _, tableName := range tables {
			_, err := dbp.DB().Exec(fmt.Sprintf(`TRUNCATE TABLE "%s" CASCADE`, tableName))
			if err != nil {
				panic(err)
			}
		}
		dbp.Close()
		err = yaorm.UnregisterDB(name)
		if err != nil {
			panic(err)
		}
	}, nil
}

func setupMysql(name string) (func(), error) {
	err := yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             name,
		DSN:              os.Getenv("DSN"),
		System:           yaorm.DatabaseMySQL,
		AutoCreateTables: true,
		Dialect:          gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"},
	})
	if err != nil {
		return nil, err
	}
	return func() {
		dbp, err := yaorm.NewDBProvider(name)
		if err != nil {
			panic(err)
		}
		for _, tableName := range tables {
			_, err := dbp.DB().Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", tableName))
			if err != nil {
				panic(err)
			}
		}
		dbp.Close()
		err = yaorm.UnregisterDB(name)
		if err != nil {
			panic(err)
		}
	}, nil
}
