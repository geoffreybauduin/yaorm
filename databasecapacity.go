package yaorm

// DatabaseCapacity describes multiple capacities provided only by some databases systems
type DatabaseCapacity int8

const (
	DatabaseCapacitySchema = iota ^ 42
)

var (
	databaseCapacities = map[DMS]map[DatabaseCapacity]bool{
		DatabaseMySQL: {
			DatabaseCapacitySchema: true,
		},
		DatabasePostgreSQL: {
			DatabaseCapacitySchema: true,
		},
		DatabaseSqlite3: {
			DatabaseCapacitySchema: false,
		},
	}
)
