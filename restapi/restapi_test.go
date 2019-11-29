package restapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/timdrysdale/dr"
	"github.com/timdrysdale/dr/mock"
)

const overlyStrict = true

func TestHandleRoot(t *testing.T) {

	// set up req & resp
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	handleRoot(resp, req)

	checkStatusCodeIs(t, resp, http.StatusNotFound)
}

func TestHandleResourcesGet(t *testing.T) {

	// set up store
	m := mock.New()
	c := map[string]int{"a": 1, "b": 2}
	m.SetCategories(c)

	// set up req & resp
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/resources", nil)
	if err != nil {
		t.Error(err.Error())
	}

	handleResourcesGet(resp, req, m)

	checkStatusCodeIs(t, resp, http.StatusOK)
	checkContentTypeContains(t, resp, "application/json")
	checkBodyEquals(t, resp, `{"a":1,"b":2}`)

}

func TestHandleResourcesGetEmptyStorage(t *testing.T) {

	// set up store
	m := mock.New()
	err := dr.ErrEmptyStorage

	m.SetError(err)

	// set up req & resp
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/resources", nil)
	if err != nil {
		t.Error(err.Error())
	}

	handleResourcesGet(resp, req, m)

	checkStatusCodeIs(t, resp, http.StatusInternalServerError)
	if overlyStrict {
		checkContentTypeContains(t, resp, "text/plain")
		checkBodyEquals(t, resp, dr.ErrEmptyStorage.Error()+"\n")
	}
}
