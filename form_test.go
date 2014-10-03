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

func TestMultipartForm(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) { testData(t, w, r) }
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	f := Form{}
	f.Field().Add("Name", "John")
	f.File().Add("File", "test-files/logo.png")

	req, err := f.Request("POST", ts.URL+"/?test=true")
	if err != nil {
		t.Error(err)
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Error(err)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		t.Error(contentType)
	}

	if _, err := f.Do("POST", ts.URL+"/?test=true"); err != nil {
		t.Error(err)
	}
}

func TestSimpleGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("Name") != "John" {
			t.Error(r.URL)
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	f := Form{}
	f.Field().Add("Name", "John")

	req, err := f.Request("GET", ts.URL)
	if err != nil {
		t.Error(err)
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Error(err)
	}

}

func TestSimplestGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query()) > 0 {
			t.Error(r.URL)
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	if _, err := (&Form{}).Do("GET", ts.URL); err != nil {
		t.Error(err)
	}
}

func TestPostForm(t *testing.T) {
	testForm("POST", t, func(r *http.Request) {
		ctype := r.Header.Get("Content-Type")
		if ctype != "application/x-www-form-urlencoded" {
			t.Fatal(ctype)
		}
		if r.Method != "POST" {
			t.Fatal(r.Method)
		}
	})
}

func TestPutForm(t *testing.T) {
	testForm("PUT", t, func(r *http.Request) {
		ctype := r.Header.Get("Content-Type")
		if ctype != "application/x-www-form-urlencoded" {
			t.Fatal(ctype)
		}
		if r.Method != "PUT" {
			t.Fatal(r.Method)
		}
	})
}

func TestGetForm(t *testing.T) {
	testForm("GET", t, func(r *http.Request) {
		if r.Method != "GET" {
			t.Fatal(r.Method)
		}
	})
}

func testForm(verb string, t *testing.T, fn func(*http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fn(r)
		testData(t, w, r)
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	f := NewForm(nil, nil)
	f.Field().Add("Name", "John")
	f.Field().Add("Likes", "Ice Cream")
	f.Header().Add("Host", "example.com")

	req, err := f.Request(verb, ts.URL+"/?test=true")
	if err != nil {
		t.Error(err)
	}

	if req.Method != verb {
		t.Log(req.Method)
		t.Fail()
	}

	if req.Header.Get("Host") != "example.com" {
		t.Fail()
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Error(err)
	}
}

func testData(t *testing.T, w http.ResponseWriter, r *http.Request) {
	ctype := r.Header.Get("Content-Type")
	m := make(url.Values)

	if strings.Contains(ctype, "json") {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&m); err != nil {
			t.Error(err)
		}
	} else if strings.Contains(ctype, "multipart") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Error(err)
		}
		m = r.Form
	} else if strings.Contains(ctype, "form") {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		m = r.Form
	} else {
		m = r.URL.Query()
	}

	if r.URL.Query().Get("test") != "true" {
		t.Error(r.URL)
	}

	if len(m["Name"]) == 0 || m["Name"][0] != "John" {
		t.Log(r.URL)
		t.Error(m)
	}

	if strings.Contains(ctype, "multipart") {
		if _, h, err := r.FormFile("File"); err != nil {
			t.Error(err)
		} else if h.Filename != "logo.png" {
			t.Error(h.Filename)
		}
	}
}

func TestFailedCase(t *testing.T) {
	f := Form{}
	f.Field().Add("Name", "John")
	f.File().Add("File", "test-files/logo1.png") //file doesn't exists

	_, err := f.Request("POST", "http://local/")
	if !os.IsNotExist(err) {
		t.Error(err)
	}
}
