package ram

import (
	"errors"

	"github.com/timdrysdale/dr"
)

type RamStorage struct {
	resources map[string]map[string]*dr.Dr
}

func (r *RamStorage) Add(dr dr.Dr) error {
	if dr.Category == "" {
		return errors.New("Undefined Category")
	}
	if dr.ID == "" {
		return errors.New("Undefined ID")
	}

	return nil
}

func (r *RamStorage) List(category string) (error, []dr.Dr) {
	drList := make([]dr.Dr, 0)
	return nil, drList
}

func (r *RamStorage) Request(category string, id string) (error, dr.Dr) {
	dr := dr.Dr{}
	return nil, dr
}

func (r *RamStorage) HealthCheck() error {
	if r.resources != nil {
		return nil
	} else {
		return errors.New("Not initialised")
	}
}

func (r *RamStorage) Reset() error {
	r.resources = make(map[string]map[string]*dr.Dr)
	return r.HealthCheck()
}

func New() dr.Storage {
	r := RamStorage{resources: make(map[string]map[string]*dr.Dr)}
	return &r
}
