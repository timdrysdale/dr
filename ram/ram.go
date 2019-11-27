package ram

import (
	"errors"
	"strings"

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

	if strings.Contains(resource.ID, dr.Separator) {
		return dr.ErrIllegalID
	}

	if strings.Contains(resource.Category, dr.Separator) {
		return dr.ErrIllegalCategory
	}

	// create category if does not already exist
	if _, ok := r.resources[resource.Category]; !ok {
		r.resources[resource.Category] = make(map[string]*dr.Dr)
	}

	r.resources[resource.Category][resource.ID] = &resource

	return nil
}

func (r *RamStorage) List(category string) (error, map[string]*dr.Dr) {

	// existence check
	if _, ok := r.resources[category]; !ok {
		return dr.ErrNoSuchCategory, nil
	}

	// empty list check
	if len(r.resources[category]) == 0 {
		return dr.ErrEmptyList, nil
	}

	publicList := make(map[string]*dr.Dr)
	// return list omitting details of the resource
	for id, resource := range r.resources[category] {
		publicList[id] = resource
		publicList[id].Resource = ""
	}

	return nil, publicList
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
