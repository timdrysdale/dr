package dr

import "errors"

type Storage interface {
	Add(dr Dr) error
	Categories() (error, map[string]int)
	Delete(category string, id string) (error, Dr)
	Get(category string, id string) (error, Dr)
	HealthCheck() error
	List(category string) (error, map[string]Dr)
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
