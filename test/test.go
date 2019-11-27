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

var addTests = []struct {
	name     string
	resource dr.Dr
	expected error
}{
	{"reject no Category or ID", dr.Dr{}, dr.ErrUndefinedCategory},
	{"reject no Id", dr.Dr{Category: "DoesNotMatter"}, dr.ErrUndefinedID},
	{"reject no Category", dr.Dr{ID: "DoesNotMatter"}, dr.ErrUndefinedCategory},
	{"reject illegal dot in ID", dr.Dr{ID: "Does.Not.Matter"}, dr.ErrIllegalID},
	{"reject illegal dot in Category", dr.Dr{Category: "Does.Not.Matter"}, dr.ErrIllegalCategory},
}

func TestInterface(t *testing.T, tester Tester) {

	// initialisation

	storage := tester.New()
	//expect blocks until ready, so test right away
	result := (storage.HealthCheck() == nil)
	processResult(t, result, "storage healthy after initialisation")

	// adding
	for _, test := range addTests {
		result = reflect.DeepEqual(storage.Add(test.resource), test.expected)
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
