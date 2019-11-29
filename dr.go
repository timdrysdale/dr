package dr

import "errors"

type Storage interface {
	Add(dr Dr) error
	Categories() (map[string]int, error)
	Delete(category string, id string) (Dr, error)
	Get(category string, id string) (Dr, error)
	HealthCheck() error
	List(category string) (map[string]Dr, error)
	Reset() error
}

type Dr struct {
	Category    string
	Description string
	ID          string
	Resource    string
	Reusable    bool
	TTL         int64
}

const Separator = "." //to ease usage of simple key-value stores, via key = <category>.<ID>

var ErrUndefinedCategory = errors.New("Undefined Category")
var ErrUndefinedID = errors.New("Undefined ID")
var ErrIllegalCategory = errors.New("Illegal Category")
var ErrIllegalID = errors.New("Illegal ID")
var ErrResourceNotFound = errors.New("Resource not found")
var ErrEmptyList = errors.New("List is empty")
var ErrEmptyStorage = errors.New("Storage is empty")
var ErrUnhealthy = errors.New("Unhealthy storage")
