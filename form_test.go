package easyreq

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func testReq(t *testing.T, w http.ResponseWriter, r *http.Request) {
	ctype := r.Header.Get("Content-Type")
	t.Log(ctype)
	m := make(url.Values)
	if strings.Contains(ctype, "json") {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&m); err != nil {
			t.Fatal(err)
		}

	} else if strings.Contains(ctype, "form") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatal(err)
		}
		m = r.Form
	} else {
		m = r.URL.Query()
	}

	if r.URL.Query().Get("test") != "true" {
		t.Fatal(r.URL.String())
	}

	if len(m["Name"]) == 0 || m["Name"][0] != "John" {
		t.Log(r.URL.String())
		t.Fatal(m)
	}

	if strings.Contains(ctype, "multipart") {
		if _, _, err := r.FormFile("File"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPostForm(t *testing.T) {
	testForm("POST", t)
}

func TestPutForm(t *testing.T) {
	testForm("PUT", t)
}

func TestMultipartForm(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) { testReq(t, w, r) }
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	f := Form{}
	f.Field().Add("Name", "John")
	f.File().Add("File", "test-files/logo.png")

	req, err := f.Request("POST", ts.URL+"/?test=true")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatal(err)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		t.Fail()
	}

	if _, err := f.Do("POST", ts.URL+"/?test=true"); err != nil {
		t.Fatal(err)
	}
}

func TestFailedCase(t *testing.T) {
	f := Form{}
	f.Field().Add("Name", "John")
	f.File().Add("File", "test-files/logo1.png") //file doesn't exists

	_, err := f.Request("POST", "http://local/")
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func TestGetForm(t *testing.T) {
	testForm("GET", t)
}

func testForm(verb string, t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) { testReq(t, w, r) }
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	f := NewForm(nil, nil)
	f.Field().Add("Name", "John")
	f.Field().Add("Likes", "Ice Cream")
	f.Header().Add("Host", "example.com")

	req, err := f.Request(verb, ts.URL+"/?test=true")
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != verb {
		t.Log(req.Method)
		t.Fail()
	}

	if req.Header.Get("Host") != "example.com" {
		t.Fail()
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatal(err)
	}
}
