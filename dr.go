package dr

import "errors"

type Storage interface {
	Add(dr Dr) error
	List(category string) (error, map[string]*Dr)
	Request(category string, id string) (error, Dr)
	HealthCheck() error
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

const Separator = "."

var ErrUndefinedCategory = errors.New("Undefined Category")
var ErrUndefinedID = errors.New("Undefined ID")
var ErrIllegalCategory = errors.New("Illegal Category")
var ErrIllegalID = errors.New("Illegal ID")
var ErrNoSuchCategory = errors.New("Category not found / does not exist")
var ErrNoSuchID = errors.New("ID not found / does not exist")
var ErrEmptyList = errors.New("List is empty")
