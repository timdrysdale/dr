package restapi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/timdrysdale/dr/mock"
)

func TestHandleGetRoot(t *testing.T) {

	// set up store
	m := mock.New()
	c := map[string]int{"a": 1, "b": 2}
	expected := `{"a":1,"b":2}`
	m.SetCategories(c)

	// set up req & resp
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	handleGetRoot(resp, req, m)

	// check
	if got, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		if strings.Contains(string(got), "Error") {
			t.Errorf("header response shouldn't return error: %s", got)
		} else if !strings.Contains(string(got), expected) {
			t.Errorf("header response doesn't match:\n%s", got)
		}
	}
}
