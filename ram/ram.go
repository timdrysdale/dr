package ram

import (
	"errors"
	"strings"
	"sync"

	"github.com/timdrysdale/dr"
)

type RamStorage struct {
	resources map[string]map[string]dr.Dr
	sync.RWMutex
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

	r.Lock()
	defer r.Unlock()

	if _, ok := r.resources[resource.Category]; !ok {
		r.resources[resource.Category] = make(map[string]dr.Dr)
	}

	r.resources[resource.Category][resource.ID] = resource

	return nil
}

func (r *RamStorage) List(category string) (error, map[string]dr.Dr) {

	publicList := make(map[string]dr.Dr)

	r.RLock()
	defer r.RUnlock()

	// existence check
	if _, ok := r.resources[category]; !ok {
		return dr.ErrNoSuchCategory, publicList
	}

	// empty list check
	if len(r.resources[category]) == 0 {
		return dr.ErrEmptyList, publicList
	}

	// return list omitting details of the resource

	for id, resource := range r.resources[category] {
		publicResource := resource
		publicResource.Resource = ""
		publicList[id] = publicResource
	}

	return nil, publicList
}

func (r *RamStorage) Get(category string, id string) (error, dr.Dr) {
	resource := dr.Dr{}

	r.Lock()
	//TODO
	// delete singleuse resource from memory - so can't pass a pointer!
	r.Unlock()

	return nil, resource
}

func (r *RamStorage) HealthCheck() error {

	r.RLock()
	defer r.RUnlock()

	if r.resources != nil {
		return nil
	} else {
		return errors.New("Not initialised")
	}
}

func (r *RamStorage) Reset() error {

	r.Lock()
	r.resources = make(map[string]map[string]dr.Dr)
	r.Unlock()

	return r.HealthCheck()
}

func (r *RamStorage) Categories() (error, []string) {
	r.RLock()
	defer r.RUnlock()
	return dr.ErrEmptyStorage, make([]string, 0)
}

func (r *RamStorage) Population() (error, map[string]int) {
	r.RLock()
	defer r.RUnlock()
	return dr.ErrEmptyStorage, make(map[string]int)
}

func New() dr.Storage {
	r := RamStorage{resources: make(map[string]map[string]dr.Dr)}
	return &r
}
