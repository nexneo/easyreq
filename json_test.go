package easyreq

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonPost(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) { testReq(t, w, r) }
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	v := struct {
		Name, Likes []string
	}{[]string{"John"}, []string{"Ice Cream"}}

	req, err := (&Json{}).Set(v).Request("POST", ts.URL)
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
