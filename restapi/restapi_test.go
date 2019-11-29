package restapi

import (
	"bytes"
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
	//multiple responses for multiple calls (#TODO but hopefully #YAGNI!)
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

func TestHandleCategoryPost(t *testing.T) {

	// set up store
	m := mock.New()

	ID1 := "some_id"
	ID2 := "other_id"

	resource1 := dr.Dr{
		Category:    "cat23",
		Description: "desc",
		ID:          ID1,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	resource2 := dr.Dr{
		Category:    "cat23",
		Description: "desc",
		ID:          ID2,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	list, err := json.Marshal(map[string]dr.Dr{ID1: resource1, ID2: resource2})
	if err != nil {
		t.Error(err)
	}
	// set up req & resp
	resp := httptest.NewRecorder()
	category := "cat23"
	r := bytes.NewReader(list)
	req, err := http.NewRequest("POST", "", r)
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"category": category,
	})

	handleCategoryPost(resp, req, m)

	if m.Method["Add"] != 2 {
		t.Errorf("Didn't call Add once, but %d times\n", m.Method["List"])
	}

	checkStatusCodeIs(t, resp, http.StatusOK)
}

func TestHandleCategoryPostCategoryError(t *testing.T) {

	// set up store
	m := mock.New()

	ID1 := "some_id"
	ID2 := "other_id"

	resource1 := dr.Dr{
		Category:    "secretCategory!", //Here's the cross end-point attack!
		Description: "desc",
		ID:          ID1,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	resource2 := dr.Dr{
		Category:    "cat23",
		Description: "desc",
		ID:          ID2,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	list, err := json.Marshal(map[string]dr.Dr{ID1: resource1, ID2: resource2})
	if err != nil {
		t.Error(err)
	}
	// set up req & resp
	resp := httptest.NewRecorder()
	category := "cat23"
	r := bytes.NewReader(list)
	req, err := http.NewRequest("POST", "", r)
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"category": category,
	})

	handleCategoryPost(resp, req, m)

	if m.Method["Add"] == 2 {
		t.Errorf("Didn't call Add once, but %d times\n", m.Method["List"])
	}
	checkStatusCodeIs(t, resp, http.StatusInternalServerError)
	checkBodyEquals(t, resp, dr.ErrIllegalCategory.Error()+":secretCategory!\n")
}

func TestHandleCategoryPostIDError(t *testing.T) {

	// set up store
	m := mock.New()

	ID1 := "some_id"
	ID2 := "other_id"

	resource1 := dr.Dr{
		Category:    "cat23",
		Description: "desc",
		ID:          ID1,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	resource2 := dr.Dr{
		Category:    "cat23",
		Description: "desc",
		ID:          ID1, //Here's the deliberate error!
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	list, err := json.Marshal(map[string]dr.Dr{ID1: resource1, ID2: resource2})
	if err != nil {
		t.Error(err)
	}
	// set up req & resp
	resp := httptest.NewRecorder()
	category := "cat23"
	r := bytes.NewReader(list)
	req, err := http.NewRequest("POST", "", r)
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"category": category,
	})

	handleCategoryPost(resp, req, m)

	if m.Method["Add"] == 2 {
		t.Errorf("Didn't call Add once, but %d times\n", m.Method["List"])
	}
	checkStatusCodeIs(t, resp, http.StatusInternalServerError)
	checkBodyEquals(t, resp, dr.ErrUndefinedID.Error()+": did you mean some_id or other_id?\n")
}

func TestHandleIDDelete(t *testing.T) {

	// set up store
	m := mock.New()
	category := "some_category"
	ID := "some_id"
	resource := dr.Dr{
		Category:    category,
		Description: "desc",
		ID:          ID,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	m.SetResource(resource)

	// set up req & resp
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "", nil)
	if err != nil {
		t.Error(err.Error())
	}
	req = mux.SetURLVars(req, map[string]string{
		"category": category,
		"id":       ID,
	})

	handleIDDelete(resp, req, m)

	if m.Method["Delete"] != 1 {
		t.Errorf("Didn't call Delete once, but %d times\n", m.Method["Delete"])
	}

	if m.GetCategory() != category {
		t.Errorf(".Delete() called with wrong category:\ngot:%s\nexp:%s\n",
			m.GetCategory(), category)
	}

	if m.GetID() != ID {
		t.Errorf(".Delete() called with wrong ID:\ngot:%s\nexp:%s\n",
			m.GetID(), ID)
	}

	checkStatusCodeIs(t, resp, http.StatusOK)
	checkBodyEquals(t, resp, "") //don't return the resource, could be trying to recover a resource issue by deleting a huge resource etc
}

func TestHandleIDGet(t *testing.T) {

	// set up store
	m := mock.New()
	category := "some_category"
	ID := "some_id"
	resource := dr.Dr{
		Category:    category,
		Description: "desc",
		ID:          ID,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}
	m.SetResource(resource)

	// set up req & resp
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Error(err.Error())
	}
	req = mux.SetURLVars(req, map[string]string{
		"category": category,
		"id":       ID,
	})

	handleIDGet(resp, req, m)

	if m.GetCategory() != category {
		t.Errorf(".List() called with wrong category:\ngot:%s\nexp:%s\n",
			m.GetCategory(), category)
	}

	if m.GetID() != ID {
		t.Errorf(".Delete() called with wrong ID:\ngot:%s\nexp:%s\n",
			m.GetID(), ID)
	}

	if m.Method["Get"] != 1 {
		t.Errorf("Didn't call Get once, but %d times\n", m.Method["Delete"])
	}

	obj, err := json.Marshal(resource)
	if err != nil {
		t.Errorf("Failed to formulate expected response")
	}
	expected := string(obj)
	checkStatusCodeIs(t, resp, http.StatusOK)
	checkContentTypeContains(t, resp, "application/json")
	checkBodyEquals(t, resp, expected)

}
func TestHandleIDPost(t *testing.T) {

	// set up store
	m := mock.New()

	ID1 := "some_id"
	category := "cat23"

	resource1 := dr.Dr{
		Category:    category,
		Description: "desc",
		ID:          ID1,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	resource, err := json.Marshal(resource1)
	if err != nil {
		t.Error(err)
	}
	// set up req & resp
	resp := httptest.NewRecorder()
	r := bytes.NewReader(resource)
	req, err := http.NewRequest("POST", "", r)
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"category": category,
		"id":       ID1,
	})

	handleIDPost(resp, req, m)

	if m.Method["Add"] != 1 {
		t.Errorf("Didn't call Add once, but %d times\n", m.Method["List"])
	}

	if m.GetResource() != resource1 {
		t.Errorf("Added resource did not match request")
	}

	checkStatusCodeIs(t, resp, http.StatusOK)
}

func TestHandleIDPostCategoryError(t *testing.T) {

	// set up store
	m := mock.New()

	ID1 := "some_id"
	category := "cat23"

	resource1 := dr.Dr{
		Category:    "secretCategory!", //attack is here!
		Description: "desc",
		ID:          ID1,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	resource, err := json.Marshal(resource1)
	if err != nil {
		t.Error(err)
	}
	// set up req & resp
	resp := httptest.NewRecorder()
	r := bytes.NewReader(resource)
	req, err := http.NewRequest("POST", "", r)
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"category": category,
		"id":       ID1,
	})

	handleIDPost(resp, req, m)

	if m.Method["Add"] != 0 {
		t.Errorf("Didn't call Add zero times, but %d times\n", m.Method["List"])
	}

	checkStatusCodeIs(t, resp, http.StatusInternalServerError)
	checkBodyEquals(t, resp, dr.ErrIllegalCategory.Error()+":secretCategory!\n")
}

func TestHandleIDPostIDError(t *testing.T) {

	// set up store
	m := mock.New()

	ID1 := "some_id"
	ID2 := "other_id"
	category := "cat23"

	resource1 := dr.Dr{
		Category:    category,
		Description: "desc",
		ID:          ID1,
		Resource:    "res",
		Reusable:    true,
		TTL:         123}

	resource, err := json.Marshal(resource1)
	if err != nil {
		t.Error(err)
	}
	// set up req & resp
	resp := httptest.NewRecorder()
	r := bytes.NewReader(resource)
	req, err := http.NewRequest("POST", "", r)
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"category": category,
		"id":       ID2, //deliberately different to resource1
	})

	handleIDPost(resp, req, m)

	if m.Method["Add"] != 0 {
		t.Errorf("Didn't call Add zero times, but %d times\n", m.Method["List"])
	}

	checkStatusCodeIs(t, resp, http.StatusInternalServerError)
	checkBodyEquals(t, resp, dr.ErrUndefinedID.Error()+": did you mean some_id or other_id?\n")
}
