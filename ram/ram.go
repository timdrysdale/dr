package ram

import (
	"errors"

	"github.com/timdrysdale/dr"
)

type RamStorage struct {
	resources map[string]map[string]*dr.Dr
}

func (r *RamStorage) Add(resource dr.Dr) error {
	if resource.Category == "" {
		return dr.ErrUndefinedCategory
	}
	if resource.ID == "" {
		return dr.ErrUndefinedID
	}

	return nil
}

func (r *RamStorage) List(category string) (error, []dr.Dr) {
	resourceList := make([]dr.Dr, 0)
	return nil, resourceList
}

func (r *RamStorage) Request(category string, id string) (error, dr.Dr) {
	resource := dr.Dr{}
	return nil, resource
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
