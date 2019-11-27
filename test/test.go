package test

import (
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

func TestInterface(t *testing.T, tester Tester) {

	// initialisation

	storage := tester.New()
	//expect blocks until ready, so test right away
	if storage.HealthCheck() != nil {
		t.Fatal("Unhealthy storage")
	}
	t.Logf("Storage Initialisation: PASS")

}

// Testing tips https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859

//https://www.toptal.com/go/your-introductory-course-to-testing-with-go

/*
A great memory-free trick for ensuring that the interface is satisfied at run time is to insert the following into our code:

var _ io.Reader = (*MockReader)(nil)

This checks the assertion but doesn’t allocate anything, which lets us make sure that the interface is correctly implemented at compile time, before the program actually runs into any functionality using it. An optional trick, but helpful.

*/
