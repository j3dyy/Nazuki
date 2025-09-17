package service

import "database/sql"

type Service struct {
}

func NewService(db *sql.DB) Service {
	return Service{}
}
