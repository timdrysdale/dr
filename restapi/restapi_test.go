package restapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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

func TestHandleResourcesDelete(t *testing.T) {

	// set up store
	m := mock.New()

	// set up req & resp
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/resources", nil)
	if err != nil {
		t.Error(err.Error())
	}

	handleResourcesGet(resp, req, m)

	checkStatusCodeIs(t, resp, http.StatusOK)

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

type Dr struct {
	Category    string
	Description string
	ID          string
	Resource    string
	Reusable    bool
	TTL         int64
}

func TestHandleCategoryGet(t *testing.T) {

	// set up store
	m := mock.New()
	resource := dr.Dr{
		Category:    "cat",
		Description: "desc",
		ID:          "id",
		Resource:    "res",
		Reusable:    true,
		TTL:         123}
	l := map[string]dr.Dr{"some_id": resource}
	m.SetList(l)

	// set up req & resp
	resp := httptest.NewRecorder()
	category := "importantcategory99"
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Error(err.Error())
	}
	req = mux.SetURLVars(req, map[string]string{
		"category": category,
	})

	handleCategoryGet(resp, req, m)

	if m.Args.Category != category {
		t.Errorf(".List() called with wrong category:\ngot:%s\nexp:%s\n",
			m.Args.Category, category)
	}

	obj, err := json.Marshal(resource)
	if err != nil {
		t.Errorf("Failed to formulate expected response")
	}
	expected := `{"some_id":` + string(obj) + "}"
	checkStatusCodeIs(t, resp, http.StatusOK)
	checkContentTypeContains(t, resp, "application/json")
	checkBodyEquals(t, resp, expected)

}
