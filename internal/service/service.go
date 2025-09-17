package service

import (
	"github.com/j3dyy/nazuki/internal/store"
)

type Service struct {
}

func NewService(store *store.Store) Service {
	return Service{}
}
