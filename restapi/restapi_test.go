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

	if m.GetCategory() != category {
		t.Errorf(".List() called with wrong category:\ngot:%s\nexp:%s\n",
			m.GetCategory(), category)
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

func TestHandleCategoryDelete(t *testing.T) {

	// set up store
	m := mock.New()
	resource := dr.Dr{
		Category:    "cat",
		Description: "desc",
		ID:          "id",
		Resource:    "res",
		Reusable:    true,
		TTL:         123}
	m.SetResource(resource)
	//let's hope API doesn't notice that we've got inconsistency between
	//list ID and resource.ID - else need to enhance mock to give
	//multiple responses for multiple calls
	ID1 := "some_id"
	ID2 := "other_id"
	l := map[string]dr.Dr{ID1: resource, ID2: resource}
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

	handleCategoryDelete(resp, req, m)

	if m.Method["List"] < 1 {
		t.Errorf("Didn't call List once, but %d times\n", m.Method["List"])
	}
	if m.Method["Delete"] != 2 {
		t.Errorf("Didn't call Delete twice, but %d times\n", m.Method["Delete"])
	}

	if m.GetCategory() != category {
		t.Errorf(".Delete() called with wrong category:\ngot:%s\nexp:%s\n",
			m.GetCategory(), category)
	}

	if m.GetID() != ID2 {
		t.Errorf(".Delete() called with wrong ID:\ngot:%s\nexp:%s\n",
			m.GetID(), ID2)
	}

	checkStatusCodeIs(t, resp, http.StatusOK)
}
