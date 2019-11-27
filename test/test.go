package test

import (
	"reflect"
	"testing"

	"github.com/timdrysdale/dr"
)

// see https://stackoverflow.com/questions/15897803/how-can-i-have-a-common-test-suite-for-multiple-packages-in-go

// functions needed for each implementation to test it
type Tester struct {
	New  func() dr.Storage
	Done func(*dr.Storage)

	// whatever you need. Leave nil if function does not apply
}

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
			Description: "Item-a.a"},
		nil},
	{"add resource a.b for list test",
		dr.Dr{
			Category:    "a",
			ID:          "b",
			Resource:    "Resource-a.b",
			Description: "Item-a.b"},
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
	{"return map of resources in known category with resource field removed",
		"a",
		nil,
		map[string]dr.Dr{
			"a": dr.Dr{
				Category:    "a",
				ID:          "a",
				Description: "Item-a.a",
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

	// add resources for list checks
	for _, test := range addForListTests {
		result = reflect.DeepEqual(storage.Add(test.resource), test.expected)
		processResult(t, result, test.name)
	}

	// list checks
	for _, test := range listTests {
		err, list := storage.List(test.category)
		result = (err == test.errExpected) && (reflect.DeepEqual(list, test.listExpected))
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
