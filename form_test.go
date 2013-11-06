package easyreq

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func testReq(t *testing.T, w http.ResponseWriter, r *http.Request) {
	ctype := r.Header.Get("Content-Type")
	t.Log(ctype)
	m := make(url.Values)
	if !strings.Contains(ctype, "json") {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		m = r.Form
	} else {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&m); err != nil {
			t.Fatal(err)
		}
	}

	if len(m["Name"]) > 0 && m["Name"][0] != "John" {
		t.Log(len(m["Name"]))
		t.Fatal(m)
	}

	if strings.Contains(ctype, "multipart") {
		if _, _, err := r.FormFile("File"); err != nil {
			t.Fatal(err)
		}
	}

	if r.ContentLength < 20 {
		t.Log(r.ContentLength)
		t.Fail()
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

	req, err := f.Request("POST", ts.URL)
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
}

func testForm(verb string, t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) { testReq(t, w, r) }
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	f := Form{}
	f.Field().Add("Name", "John")
	f.Field().Add("Likes", "Ice Cream")

	req, err := f.Request(verb, ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatal(err)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/x-www-form-urlencoded" {
		t.Fail()
	}
}
