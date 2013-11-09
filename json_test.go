package easyreq

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func payload() interface{} {
	v := struct {
		Name, Likes []string
	}{[]string{"John"}, []string{"Ice Cream"}}

	return v
}

func TestJsonPost(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) { testReq(t, w, r) }
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	j := Json{}
	j.Header().Add("Host", "example.com")
	req, err := j.Set(payload()).Request("POST", ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Fail()
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatal(err)
	}

	if req.Header.Get("Host") != "example.com" {
		t.Fail()
	}

	if _, err := Do(&j, "POST", ts.URL); err != nil {
		t.Fatal(err)
	}
}

func TestJsonPut(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) { testReq(t, w, r) }
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	req, err := NewJson(payload()).Request("PUT", ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Fail()
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatal(err)
	}
}
