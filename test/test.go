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
	New  func() dr.Storage
	Done func(*dr.Storage)

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
		dr.ErrResourceNotFound,
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
var expectedCategoriesAfterTTL = map[string]int{"a": 1, "x": 1}

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
			TTL:         5,
		},
		nil},
	{"add resource a.e for ttl test",
		dr.Dr{
			Category:    "a",
			ID:          "e",
			Resource:    "Resource-a.e",
			Description: "Item-a.e",
			TTL:         3,
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
	{"list after 0.001 sec shows a.a, a.c and a.d",
		1 * time.Millisecond,
		"a",
		nil,
		map[string]dr.Dr{
			"a": dr.Dr{
				Category:    "a",
				ID:          "a",
				Description: "Item-a.a",
				Reusable:    true,
			},
			"c": dr.Dr{
				Category:    "a",
				ID:          "c",
				Description: "Item-a.c",
				TTL:         1,
			},
			"d": dr.Dr{
				Category:    "a",
				ID:          "d",
				Description: "Item-a.d",
				TTL:         5,
			},
			"e": dr.Dr{
				Category:    "a",
				ID:          "e",
				Description: "Item-a.e",
				TTL:         3,
			},
		}},
	{"list after 2.001 sec shows a.a & a.d + a.e with updated TTL",
		2000 * time.Millisecond,
		"a",
		nil,
		map[string]dr.Dr{
			"a": dr.Dr{
				Category:    "a",
				ID:          "a",
				Description: "Item-a.a",
				Reusable:    true,
			},
			"d": dr.Dr{
				Category:    "a",
				ID:          "d",
				Description: "Item-a.d",
				TTL:         3,
			},
			"e": dr.Dr{
				Category:    "a",
				ID:          "e",
				Description: "Item-a.e",
				TTL:         1,
			},
		}},
}

var postTTLGetTests = []struct {
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
	{"get resource a.e",
		"a",
		"e",
		dr.ErrResourceNotFound,
		dr.Dr{},
	},
}

var postTTLDeleteTests = []struct {
	name             string
	category         string
	ID               string
	errExpected      error
	resourceExpected dr.Dr
}{
	{"delete reusable resource a.a",
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
	{"throw err deleting deleted resource a.a",
		"a",
		"a",
		dr.ErrResourceNotFound, // a.d has gone by now so category error, not ID error
		dr.Dr{},
	},
}

var listAfterTTLTests = []struct {
	name         string
	category     string
	errExpected  error
	listExpected map[string]dr.Dr
}{
	{"list after TTL tests shows a.a & a.d with updated TTL",
		"a",
		nil,
		map[string]dr.Dr{
			"a": dr.Dr{
				Category:    "a",
				ID:          "a",
				Description: "Item-a.a",
				Reusable:    true,
			},
			"d": dr.Dr{
				Category:    "a",
				ID:          "d",
				Description: "Item-a.d",
				TTL:         1,
			},
		}},
}

func TestInterface(t *testing.T, tester Tester) {

	// initialisation
	storage := tester.New() // expect New() blocks until initialisation complete
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
	categories, err := storage.Categories()
	result = (err == dr.ErrEmptyStorage) && reflect.DeepEqual(categories, map[string]int{})
	processResult(t, result, "map categories throws ErrEmptyStorage when store empty")

	// add resources for list checks
	for _, test := range addForListTests {
		result = reflect.DeepEqual(storage.Add(test.resource), test.expected)
		processResult(t, result, test.name)
	}

	// list tests
	for _, test := range listTests {
		list, err := storage.List(test.category)
		result = (err == test.errExpected) && (reflect.DeepEqual(list, test.listExpected))
		if debugTest {
			t.Log(list)
			t.Log(test.listExpected)
			t.Log(reflect.DeepEqual(list, test.listExpected))
		}
		processResult(t, result, test.name)
	}

	// categories test
	categories, err = storage.Categories()
	result = (err == nil) && reflect.DeepEqual(categories, expectedCategories)
	processResult(t, result, "map categories and number of items therein")

	// get tests
	for _, test := range getTests {
		resource, err := storage.Get(test.category, test.ID)
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
		list, err := storage.List(test.category)
		result = (err == test.errExpected) && (reflect.DeepEqual(list, test.listExpected))
		if debugTest {
			t.Log(list)
			t.Log(test.listExpected)
			t.Log(reflect.DeepEqual(list, test.listExpected))
		}
		processResult(t, result, test.name)
	}

	// Tarantino time: tests to come after TTL testing - defer so we don't skip
	defer func() {

		for _, test := range postTTLDeleteTests {
			resource, err := storage.Delete(test.category, test.ID)
			result = (err == test.errExpected) && (reflect.DeepEqual(resource, test.resourceExpected))
			if debugTest {
				t.Log(err)
				t.Log(test.errExpected)
				t.Log(resource)
				t.Log(test.resourceExpected)
				t.Log(reflect.DeepEqual(resource, test.resourceExpected))
			}
			processResult(t, result, test.name)
		}

	}()

	// nothing after this point will run if the test is -short
	if testing.Short() {
		t.Skip("**SKIP** skipping TTL tests - check before releasing though!")
	}

	// add resources for TTL test
	for _, test := range addForTTLTests {
		result = reflect.DeepEqual(storage.Add(test.resource), test.expected)
		processResult(t, result, test.name)
	}

	// post-get list tests
	for _, test := range listForTTLTests {
		time.Sleep(test.duration)
		list, err := storage.List(test.category)
		result = (err == test.errExpected) && (reflect.DeepEqual(list, test.listExpected))
		if debugTest {
			t.Log(list)
			t.Log(test.listExpected)
			t.Log(reflect.DeepEqual(list, test.listExpected))
		}
		processResult(t, result, test.name)
	}

	// await a.e expiring since last list, to ensure Get() is checking its stale
	time.Sleep(2000 * time.Millisecond)
	for _, test := range postTTLGetTests {
		resource, err := storage.Get(test.category, test.ID)
		result = (err == test.errExpected) && (reflect.DeepEqual(resource, test.resourceExpected))
		if debugTest {
			t.Log(err)
			t.Log(test.errExpected)
			t.Log(resource)
			t.Log(test.resourceExpected)
			t.Log(reflect.DeepEqual(resource, test.resourceExpected))
		}
		processResult(t, result, test.name)
	}

	// post-get list tests
	for _, test := range listAfterTTLTests {
		list, err := storage.List(test.category)
		result = (err == test.errExpected) && (reflect.DeepEqual(list, test.listExpected))
		if debugTest {
			t.Log(list)
			t.Log(test.listExpected)
			t.Log(reflect.DeepEqual(list, test.listExpected))
		}
		processResult(t, result, test.name)
	}

	// await a.d expiring to check that categories cleans stale entries
	time.Sleep(2000 * time.Millisecond)
	categories, err = storage.Categories()
	result = (err == nil) && reflect.DeepEqual(categories, expectedCategoriesAfterTTL)
	processResult(t, result, "map categories and ignore stale resource")

}

func processResult(t *testing.T, result bool, name string) {
	if result {
		t.Logf("  pass   %s\n", name)
	} else {
		t.Errorf("**FAIL** %s\n", name)
	}
}
