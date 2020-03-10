package yaorm

// DBSpecific is an interface describing how specificities of each DMS
// are handled in yaorm
type DBSpecific interface {
	// OnSessionCreated is a function executed once per DBProvider instance
	OnSessionCreated(DBProvider) error
}

type noopSpecific struct{}

func (noopSpecific) OnSessionCreated(_ DBProvider) error { return nil }
