package ram

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/timdrysdale/dr"
)

/*

Mock time - where to store the clock, and how to sub a real clock in production?
use the New constructor for this ...

https://github.com/jonboulle/clockwork

or.... is this an implementation detail that should not spill into the interface?

*/

type expiringResource struct {
	resource   dr.Dr
	validUntil int64
}

type RamStorage struct {
	resources map[string]map[string]expiringResource
	clock     clockwork.Clock
	sync.RWMutex
}

func (r *RamStorage) Now() int64 {
	return time.Now().Unix()
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
		r.resources[resource.Category] = make(map[string]expiringResource)
	}

	validUntil := int64(0) //default to living forever

	if resource.TTL > 0 {
		validUntil = r.Now() + resource.TTL
	}

	r.resources[resource.Category][resource.ID] = expiringResource{resource: resource, validUntil: validUntil}

	return nil
}

func (r *RamStorage) List(category string) (error, map[string]dr.Dr) {

	publicList := make(map[string]dr.Dr)

	r.Lock() //need a write lock because we might clean stale entries
	defer r.Unlock()

	// existence check
	if _, ok := r.resources[category]; !ok {
		return dr.ErrNoSuchCategory, publicList
	}

	// empty list check
	if len(r.resources[category]) == 0 {
		return dr.ErrEmptyList, publicList
	}

	// return list omitting details of the resource

	for id, expiringResource := range r.resources[category] {

		//clean stale entries

		expired := false

		if expiringResource.validUntil > 0 {
			// expirable
			newTTL := expiringResource.validUntil - r.Now()
			if newTTL < 0 {
				expired = true
				delete(r.resources[category], id)
			} else {
				// update TTL
				temp := r.resources[category][id]
				temp.resource.TTL = newTTL
				r.resources[category][id] = temp
			}

		}

		if !expired {
			publicResource := r.resources[category][id].resource
			publicResource.Resource = ""
			publicList[id] = publicResource
		}

	}

	return nil, publicList
}

func (r *RamStorage) Get(category string, id string) (error, dr.Dr) {

	emptyResource := dr.Dr{}

	r.Lock()
	defer r.Unlock()

	// category existence check
	if _, ok := r.resources[category]; !ok {
		return dr.ErrNoSuchCategory, emptyResource
	}

	// ID existence check
	if resource, ok := r.resources[category][id]; ok {

		if !resource.resource.Reusable {
			delete(r.resources[category], id)
		}

		return nil, resource.resource

	} else {
		return dr.ErrNoSuchID, emptyResource
	}

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
	r.resources = make(map[string]map[string]expiringResource)
	r.Unlock()

	return r.HealthCheck()
}

func (r *RamStorage) Categories() (error, map[string]int) {

	categoryMap := make(map[string]int)

	r.RLock()
	defer r.RUnlock()

	if len(r.resources) == 0 {
		return dr.ErrEmptyStorage, categoryMap
	}

	for category, resourceMap := range r.resources {
		categoryMap[category] = len(resourceMap)
	}

	return nil, categoryMap
}

func New() dr.Storage {
	r := RamStorage{resources: make(map[string]map[string]expiringResource)}
	return &r
}
