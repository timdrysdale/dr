package restapi

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

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
