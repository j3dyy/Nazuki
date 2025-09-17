package store

import "database/sql"

type Store struct {
}

func NewStore(db *sql.DB) *Store {
	return &Store{}
}
