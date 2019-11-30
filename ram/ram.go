package ram

import (
	"strings"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/timdrysdale/dr"
)

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

func (r *RamStorage) Categories() (map[string]int, error) {

	categoryMap := make(map[string]int)
	categoryList := []string{}

	r.RLock()
	numCategories := len(r.resources)
	r.RUnlock()

	if numCategories == 0 {
		return categoryMap, dr.ErrEmptyStorage
	}

	r.RLock()
	for category, _ := range r.resources {
		categoryList = append(categoryList, category)
	}
	r.RUnlock()

	for _, category := range categoryList {
		_, _ = r.List(category) //use list to do stale cleaning
	}

	r.RLock()
	numCategories = len(r.resources)
	r.RUnlock()

	if numCategories == 0 {
		return categoryMap, dr.ErrEmptyStorage
	}

	r.RLock()
	for category, resourceMap := range r.resources {
		categoryMap[category] = len(resourceMap)
	}
	r.RUnlock()

	return categoryMap, nil
}

func (r *RamStorage) Delete(category string, id string) (dr.Dr, error) {

	emptyResource := dr.Dr{}

	r.Lock()
	defer r.Unlock()

	// category existence check
	if _, ok := r.resources[category]; !ok {
		return emptyResource, dr.ErrResourceNotFound
	}

	// ID existence check & deletion
	if expiringResource, ok := r.resources[category][id]; ok {
		delete(r.resources[category], id)
		if len(r.resources[category]) == 0 {
			delete(r.resources, category)
		}
		return expiringResource.resource, nil
	} else {
		// not found
		return emptyResource, dr.ErrResourceNotFound
	}

}

func (r *RamStorage) Get(category string, id string) (dr.Dr, error) {

	emptyResource := dr.Dr{}

	r.Lock()
	defer r.Unlock()

	// category existence check
	if _, ok := r.resources[category]; !ok {

		return emptyResource, dr.ErrResourceNotFound

	}

	// ID existence check
	if expiringResource, ok := r.resources[category][id]; ok {

		//clean stale entry if found

		expired := false

		if expiringResource.validUntil > 0 {
			// expirable
			newTTL := expiringResource.validUntil - r.Now()
			if newTTL < 0 {
				expired = true
				delete(r.resources[category], id)
				if len(r.resources[category]) == 0 {
					delete(r.resources, category)
				}
			} else {
				// update TTL
				temp := r.resources[category][id]
				temp.resource.TTL = newTTL
				r.resources[category][id] = temp
				expiringResource.resource.TTL = newTTL
			}

		}

		if expired {

			// expired since last clean, don't return it

			return emptyResource, dr.ErrResourceNotFound

		} else {

			// delete if single use
			if !expiringResource.resource.Reusable {
				delete(r.resources[category], id)
				if len(r.resources[category]) == 0 {
					delete(r.resources, category)
				}

			}

			// return resource (with up-to-date TTL)
			return expiringResource.resource, nil

		}

	} else {
		//not found
		return emptyResource, dr.ErrResourceNotFound
	}

}

func (r *RamStorage) HealthCheck() error {

	r.RLock()
	defer r.RUnlock()

	if r.resources != nil {
		return nil
	} else {
		return dr.ErrUnhealthy
	}
}

func (r *RamStorage) List(category string) (map[string]dr.Dr, error) {

	publicList := make(map[string]dr.Dr)

	r.Lock() //need a write lock because we might clean stale entries
	defer r.Unlock()

	// existence check
	if _, ok := r.resources[category]; !ok {
		return publicList, dr.ErrResourceNotFound
	}

	// empty list check
	if len(r.resources[category]) == 0 {
		return publicList, dr.ErrEmptyList
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
				if len(r.resources[category]) == 0 {
					delete(r.resources, category)
				}
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

	return publicList, nil
}

func New() dr.Storage {
	r := RamStorage{resources: make(map[string]map[string]expiringResource)}
	return &r
}

func (r *RamStorage) Reset() error {

	r.Lock()
	r.resources = make(map[string]map[string]expiringResource)
	r.Unlock()

	return r.HealthCheck()
}
