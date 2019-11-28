package test

import (
	"reflect"
	"testing"
	"time"

	"github.com/timdrysdale/dr"
)

// see https://stackoverflow.com/questions/15897803/how-can-i-have-a-common-test-suite-for-multiple-packages-in-go

// functions needed for each implementation to test it
type Tester struct {
	NewForTest func() dr.Storage
	Done       func(*dr.Storage)

	// whatever you need. Leave nil if function does not apply
}

var debugTest = false

var addSanityTests = []struct {
	name     string
	resource dr.Dr
	expected error
}{
	{"reject no Category or ID", dr.Dr{}, dr.ErrUndefinedCategory},
	{"reject no Id", dr.Dr{Category: "DoesNotMatter"}, dr.ErrUndefinedID},
	{"reject no Category", dr.Dr{ID: "DoesNotMatter"}, dr.ErrUndefinedCategory},
	{"reject illegal dot in ID", dr.Dr{Category: "a", ID: "Does.Not.Matter"}, dr.ErrIllegalID},
	{"reject illegal dot in Category", dr.Dr{Category: "Does.Not.Matter", ID: "a"}, dr.ErrIllegalCategory},
	{"accept resource with nil resource, description, ttl", dr.Dr{Category: "a", ID: "0"}, nil},
	{"accept resource with zero ttl", dr.Dr{Category: "a", ID: "0", TTL: 0}, nil},
}

var addForListTests = []struct {
	name     string
	resource dr.Dr
	expected error
}{
	{"add resource a.a for list test",
		dr.Dr{
			Category:    "a",
			ID:          "a",
			Resource:    "Resource-a.a",
			Description: "Item-a.a",
			Reusable:    true},
		nil},
	{"add resource a.b for list test",
		dr.Dr{
			Category:    "a",
			ID:          "b",
			Resource:    "Resource-a.b",
			Description: "Item-a.b"},
		nil},
	{"add resource x.y for list test",
		dr.Dr{
			Category:    "x",
			ID:          "y",
			Resource:    "Resource-x.y",
			Description: "Item-x.y"},
		nil},
}

var listTests = []struct {
	name         string
	category     string
	errExpected  error
	listExpected map[string]dr.Dr
}{
	{"throw error on listing nonexistent category",
		"foo",
		dr.ErrNoSuchCategory,
		make(map[string]dr.Dr),
	},
	{"return map-by-id of one resource in category 'x' with resource field removed",
		"x",
		nil,
		map[string]dr.Dr{
			"y": dr.Dr{
				Category:    "x",
				ID:          "y",
				Description: "Item-x.y",
			},
		}},
	{"return map-by-id of two resources in category 'a' with resource field removed",
		"a",
		nil,
		map[string]dr.Dr{
			"a": dr.Dr{
				Category:    "a",
				ID:          "a",
				Description: "Item-a.a",
				Reusable:    true,
			},
			"b": dr.Dr{
				Category:    "a",
				ID:          "b",
				Description: "Item-a.b",
			},
		}},
}

var expectedCategories = map[string]int{"a": 2, "x": 1}

var getTests = []struct {
	name             string
	category         string
	ID               string
	errExpected      error
	resourceExpected dr.Dr
}{
	{"get resource a.a",
		"a",
		"a",
		nil,
		dr.Dr{
			Category:    "a",
			ID:          "a",
			Description: "Item-a.a",
			Resource:    "Resource-a.a",
			Reusable:    true,
		},
	},
	{"get resource a.b",
		"a",
		"b",
		nil,
		dr.Dr{
			Category:    "a",
			ID:          "b",
			Description: "Item-a.b",
			Resource:    "Resource-a.b",
		},
	},
}

var postGetListTests = []struct {
	name         string
	category     string
	errExpected  error
	listExpected map[string]dr.Dr
}{
	{"return map-by-id of remaining resource in category 'a'",
		"a",
		nil,
		map[string]dr.Dr{
			"a": dr.Dr{
				Category:    "a",
				ID:          "a",
				Description: "Item-a.a",
				Reusable:    true,
			},
		}},
}

var addForTTLTests = []struct {
	name     string
	resource dr.Dr
	expected error
}{
	{"add resource a.c for ttl test",
		dr.Dr{
			Category:    "a",
			ID:          "c",
			Resource:    "Resource-a.c",
			Description: "Item-a.c",
			TTL:         1,
		},
		nil},
	{"add resource a.d for ttl test",
		dr.Dr{
			Category:    "a",
			ID:          "d",
			Resource:    "Resource-a.d",
			Description: "Item-a.d",
			TTL:         2,
		},
		nil},
}

var listForTTLTests = []struct {
	name         string
	duration     time.Duration
	category     string
	errExpected  error
	listExpected map[string]dr.Dr
}{
	{"list after 0.5 sec shows a.a, a.b, a.c, a.d",
		500 * time.Millisecond,
		"a",
		nil,
		map[string]dr.Dr{
			"a": dr.Dr{
				Category:    "a",
				ID:          "a",
				Description: "Item-a.a",
				Reusable:    true,
			},
			"b": dr.Dr{
				Category:    "a",
				ID:          "b",
				Description: "Item-a.b",
			},
			"c": dr.Dr{
				Category:    "a",
				ID:          "c",
				Description: "Item-a.c",
			},
			"d": dr.Dr{
				Category:    "a",
				ID:          "d",
				Description: "Item-a.d",
			},
		}},
	{"list after 1.5 sec shows a.a, a.b, a.d",
		1500 * time.Millisecond,
		"a",
		nil,
		map[string]dr.Dr{
			"a": dr.Dr{
				Category:    "a",
				ID:          "a",
				Description: "Item-a.a",
				Reusable:    true,
			},
			"b": dr.Dr{
				Category:    "a",
				ID:          "b",
				Description: "Item-a.b",
			},
			"d": dr.Dr{
				Category:    "a",
				ID:          "d",
				Description: "Item-a.d",
			},
		}},
	{"list after 2.5 sec shows a.a, a.b",
		2500 * time.Millisecond,
		"a",
		nil,
		map[string]dr.Dr{
			"a": dr.Dr{
				Category:    "a",
				ID:          "a",
				Description: "Item-a.a",
				Reusable:    true,
			},
			"b": dr.Dr{
				Category:    "a",
				ID:          "b",
				Description: "Item-a.b",
			},
		}},
}

func TestInterface(t *testing.T, tester Tester) {

	// initialisation
	storage := tester.NewForTest() // expect New() blocks until initialisation complete
	result := (storage.HealthCheck() == nil)
	processResult(t, result, "storage healthy after initialisation")

	// add - sanity checks
	for _, test := range addSanityTests {
		result = reflect.DeepEqual(storage.Add(test.resource), test.expected)
		processResult(t, result, test.name)
	}

	// reset
	result = (storage.Reset() == nil)
	processResult(t, result, "storage healthy after reset")

	// categories on empty storage
	err, categories := storage.Categories()
	result = (err == dr.ErrEmptyStorage) && reflect.DeepEqual(categories, map[string]int{})
	processResult(t, result, "map categories throws ErrEmptyStorage when store empty")

	// add resources for list checks
	for _, test := range addForListTests {
		result = reflect.DeepEqual(storage.Add(test.resource), test.expected)
		processResult(t, result, test.name)
	}

	// list tests
	for _, test := range listTests {
		err, list := storage.List(test.category)
		result = (err == test.errExpected) && (reflect.DeepEqual(list, test.listExpected))
		if debugTest {
			t.Log(list)
			t.Log(test.listExpected)
			t.Log(reflect.DeepEqual(list, test.listExpected))
		}
		processResult(t, result, test.name)
	}

	// categories test
	err, categories = storage.Categories()
	result = (err == nil) && reflect.DeepEqual(categories, expectedCategories)
	processResult(t, result, "map categories and number of items therein")

	// get tests
	for _, test := range getTests {
		err, resource := storage.Get(test.category, test.ID)
		result = (err == test.errExpected) && (reflect.DeepEqual(resource, test.resourceExpected))
		if debugTest {
			t.Log(resource)
			t.Log(test.resourceExpected)
			t.Log(reflect.DeepEqual(resource, test.resourceExpected))
		}
		processResult(t, result, test.name)
	}

	// post-get list tests
	for _, test := range postGetListTests {
		err, list := storage.List(test.category)
		result = (err == test.errExpected) && (reflect.DeepEqual(list, test.listExpected))
		if debugTest {
			t.Log(list)
			t.Log(test.listExpected)
			t.Log(reflect.DeepEqual(list, test.listExpected))
		}
		processResult(t, result, test.name)
	}

	// add resources for TTL test
	for _, test := range addForTTLTests {
		result = reflect.DeepEqual(storage.Add(test.resource), test.expected)
		processResult(t, result, test.name)
	}

	// post-get list tests
	for _, test := range listForTTLTests {
		time.Sleep(test.duration)
		err, list := storage.List(test.category)
		result = (err == test.errExpected) && (reflect.DeepEqual(list, test.listExpected))
		if debugTest {
			t.Log(list)
			t.Log(test.listExpected)
			t.Log(reflect.DeepEqual(list, test.listExpected))
		}
		processResult(t, result, test.name)
	}
}

func processResult(t *testing.T, result bool, name string) {
	if result {
		t.Logf("  pass   %s\n", name)
	} else {
		t.Errorf("**FAIL** %s\n", name)
	}
}

// Testing tips https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859

//https://www.toptal.com/go/your-introductory-course-to-testing-with-go

/*
A great memory-free trick for ensuring that the interface is satisfied at run time is to insert the following into our code:

var _ io.Reader = (*MockReader)(nil)

This checks the assertion but doesn’t allocate anything, which lets us make sure that the interface is correctly implemented at compile time, before the program actually runs into any functionality using it. An optional trick, but helpful.

*/
