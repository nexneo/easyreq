package easyreq

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Json struct {
	payload interface{}
	header  http.Header
}

func NewJson(payload interface{}) *Json {
	j := &Json{}
	j.payload = payload
	return j
}

func (j *Json) Set(payload interface{}) *Json {
	j.payload = payload
	return j
}

func (j *Json) Header() http.Header {
	if j.header == nil {
		j.header = make(http.Header)
	}
	return j.header
}

// Helper funcation send requests using http.DefaultClient
func (j *Json) Do(verb, urlStr string) (*http.Response, error) {
	return do(j, verb, urlStr)
}

func (j *Json) SetBasicAuth(username, password string) {
	setBasicAuth(j, username, password)
}

// Returns request based on current payload assoicated with Json request,
// it will always have correct Content-Type set
func (j *Json) Request(verb, urlStr string) (req *http.Request, err error) {
	var data []byte

	if j.payload != nil {
		data, err = json.Marshal(j.payload)
		if err != nil {
			return
		}
	}

	body := bytes.NewBuffer(data)
	req, err = http.NewRequest(verb, urlStr, body)
	if err != nil {
		return
	}

	req.Header = j.Header()
	req.Header.Set("Content-Type", "application/json")
	return
}
