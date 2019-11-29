// package mock to allow testing APIs without using an actual dr.Storage implementation
package mock

import (
	"github.com/timdrysdale/dr"
)

type In struct {
	Category string
	ID       string
	Resource dr.Dr
}

type Out struct {
	Error      error
	Categories map[string]int
	Resource   dr.Dr
	List       map[string]dr.Dr
}

type MockStorage struct {
	Args    In
	Returns Out
	Method  map[string]int
}

// instantiation

func New() *MockStorage {
	return &MockStorage{}
}

// mock methods for setting return values

func (m *MockStorage) SetCategories(c map[string]int) {
	m.Returns.Categories = c
}

func (m *MockStorage) SetList(l map[string]dr.Dr) {
	m.Returns.List = l
}

func (m *MockStorage) SetError(err error) {
	m.Returns.Error = err
}

func (m *MockStorage) SetResource(r dr.Dr) {
	m.Returns.Resource = r
}

// mock methods for getting arguments supplied
func (m *MockStorage) GetAdd() dr.Dr {
	return m.Args.Resource
}

func (m *MockStorage) GetCategory() string {
	return m.Args.Category
}

func (m *MockStorage) GetID() string {
	return m.Args.ID
}

func (m *MockStorage) GetMethod() map[string]int {
	return m.Method
}

// method for updating call record

func (m *MockStorage) logMethod(method string) {
	if m.Method == nil {
		m.Method = make(map[string]int)
	}
	if _, ok := m.Method[method]; ok {
		m.Method[method] += 1
	} else {
		m.Method[method] = 1
	}
}

// interface methods
func (m *MockStorage) Add(resource dr.Dr) error {
	m.logMethod("Add")
	m.Args.Resource = resource
	return m.Returns.Error
}

func (m *MockStorage) Categories() (map[string]int, error) {
	m.logMethod("Categories")
	return m.Returns.Categories, m.Returns.Error
}

func (m *MockStorage) Delete(category string, id string) (dr.Dr, error) {
	m.logMethod("Delete")
	m.Args.Category = category
	m.Args.ID = id
	return m.Returns.Resource, m.Returns.Error
}

func (m *MockStorage) Get(category string, id string) (dr.Dr, error) {
	m.logMethod("Get")
	m.Args.Category = category
	m.Args.ID = id
	return m.Returns.Resource, m.Returns.Error
}

func (m *MockStorage) HealthCheck() error {
	m.logMethod("HealthCheck")
	return m.Returns.Error
}

func (m *MockStorage) List(category string) (map[string]dr.Dr, error) {
	m.logMethod("List")
	m.Args.Category = category
	return m.Returns.List, m.Returns.Error
}

func (m *MockStorage) Reset() error {
	m.logMethod("Reset")
	return m.Returns.Error
}
