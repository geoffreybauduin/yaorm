package yaorm

// PostgresSpecific holds specific configurations used by Postgres
type PostgresSpecific struct {
	// IntervaStyle holds the style of output of the interval
	// See Postgres manual 8.5.5
	IntervalStyle string
}

// OnSessionCreated is executed when a Session is created.
// This function will use the:
// - IntervalStyle
func (p PostgresSpecific) OnSessionCreated(dbp DBProvider) error {
	if p.IntervalStyle != "" {
		if _, err := dbp.DB().Exec("SET intervalstyle = ?", p.IntervalStyle); err != nil {
			return err
		}
	}
	return nil
}
