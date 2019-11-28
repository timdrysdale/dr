package ram

import (
	"testing"

	"github.com/timdrysdale/dr/test"
)

// run generic tests on this particular implementation
func TestInterface(t *testing.T) {
	t.Log("Testing ./ram ...")
	test.TestInterface(t, test.Tester{NewForTest: NewForTest})
}
