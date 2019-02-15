package models

type DB struct {
	// TODO: Store users in a database
	Db map[string]*User
}

func NewDB(dataSourceName string) (*DB, error) {
	return &DB{make(map[string]*User)}, nil
}
