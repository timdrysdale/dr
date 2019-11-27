package dr

import "errors"

type Storage interface {
	Add(dr Dr) error
	List(category string) (error, []Dr)
	Request(category string, id string) (error, Dr)
	HealthCheck() error
}

type Dr struct {
	Category    string
	Description string
	ID          string
	Resource    string
	Reusable    bool
	TTL         int64
}

var ErrUndefinedCategory = errors.New("Undefined Category")
var ErrUndefinedID = errors.New("Undefined ID")
var ErrIllegalCategory = errors.New("Illegal Category")
var ErrIllegalID = errors.New("Illegal ID")
