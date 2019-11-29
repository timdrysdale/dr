package restapi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/timdrysdale/dr"
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

	checkStatusCodeIs(t, resp, 200)
	checkContentTypeContains(t, resp, "application/json")
	checkBodyEquals(t, resp, expected)

}

func TestHandleGetRootEmptyStorage(t *testing.T) {

	// set up store
	m := mock.New()
	err := dr.ErrEmptyStorage
	expected := dr.ErrEmptyStorage.Error() + "\n"
	m.SetError(err)

	// set up req & resp
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	handleGetRoot(resp, req, m)

	checkStatusCodeIs(t, resp, 500)
	checkContentTypeContains(t, resp, "text/plain")
	checkBodyEquals(t, resp, expected)
}

func checkStatusCodeIs(t *testing.T, resp *httptest.ResponseRecorder, expected int) {
	got := (resp.Result()).StatusCode
	if got != expected {
		t.Errorf("Unexpected StatusCode:\ngot:%v\nexp:%v\n", got, expected)
	}
}

func checkContentTypeContains(t *testing.T, resp *httptest.ResponseRecorder, expected string) {
	ct := (resp.Header())["Content-Type"]
	found := false
	for _, ctype := range ct {
		if strings.Contains(ctype, expected) {
			found = true
		}
	}
	if !found {
		t.Errorf("\nError:Unexpected Content-type:\nexp:[%s]\ngot:%v\n", expected, ct)
	}
}

func checkBodyEquals(t *testing.T, resp *httptest.ResponseRecorder, expected string) {

	if got, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err.Error())
	} else {
		if string(got) != expected {
			t.Errorf("\nError:Unexpected response.Body:\ngot:%v\nexp:%v\n", strings.TrimSuffix(string(got), "\n"), expected)
		}
	}
}
