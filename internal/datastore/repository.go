package datastore

import (
	"errors"

	"github.com/flohansen/dasher-server/internal/model"
)

type InMemDatastore struct {
}

func NewInMem() *InMemDatastore {
	return &InMemDatastore{}
}

func (repo *InMemDatastore) GetAll() ([]model.FeatureData, error) {
	return nil, errors.New("not implemented")
}
