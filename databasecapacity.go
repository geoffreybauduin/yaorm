package yaorm

// DatabaseCapacity describes multiple capacities provided only by some databases systems
type DatabaseCapacity int8

const (
	DatabaseCapacitySchema = iota ^ 42
	DatabaseCapacityUUID
)

var (
	databaseCapacities = map[DMS]map[DatabaseCapacity]bool{
		DatabaseMySQL: {
			DatabaseCapacitySchema: true,
			DatabaseCapacityUUID:   false,
		},
		DatabasePostgreSQL: {
			DatabaseCapacitySchema: true,
			DatabaseCapacityUUID:   true,
		},
		DatabaseSqlite3: {
			DatabaseCapacitySchema: false,
			DatabaseCapacityUUID:   false,
		},
	}
)
