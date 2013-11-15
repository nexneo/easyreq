package easyreq

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Form struct {
	fields url.Values
	files  url.Values
	header http.Header
}

type Requester interface {
	Header() http.Header
	Request(string, string) (*http.Request, error)
}

// Creates new Form with given fields and files
func NewForm(fields, files url.Values) *Form {
	f := new(Form)
	if fields != nil {
		f.fields = fields
	}
	if files != nil {
		f.files = files
	}
	return f
}

// Returns url.Values which should be used to Add form field
func (f *Form) Field() url.Values {
	if f.fields == nil {
		f.fields = make(url.Values)
	}
	return f.fields
}

// Returns url.Values which should be used to Add form file
func (f *Form) File() url.Values {
	if f.files == nil {
		f.files = make(url.Values)
	}
	return f.files
}

func (f *Form) Header() http.Header {
	if f.header == nil {
		f.header = make(http.Header)
	}
	return f.header
}

func setBasicAuth(r Requester, username, password string) {
	s := username + ":" + password
	r.Header().Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s)))
}

func (f *Form) SetBasicAuth(username, password string) {
	setBasicAuth(f, username, password)
}

func do(r Requester, verb, urlStr string) (*http.Response, error) {
	req, err := r.Request(verb, urlStr)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (f *Form) Do(verb, urlStr string) (*http.Response, error) {
	return do(f, verb, urlStr)
}

// Returns request based on current Fields and Files assoicated with form
// Request will always have correct Content-Type set
func (form *Form) Request(verb, urlStr string) (*http.Request, error) {
	if verb == "GET" {
		return getRequest(form, verb, urlStr)
	}

	if len(form.files) == 0 {
		return postWithoutFiles(form, verb, urlStr)
	}

	body := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(body)

	for key, paths := range form.files {
		for _, pathVal := range paths {
			pathVal = filepath.Clean(pathVal)
			file, err := os.Open(pathVal)
			if err != nil {
				return nil, err
			}

			mimePart, err := bodyWriter.CreateFormFile(key, filepath.Base(pathVal))
			if err != nil {
				return nil, err
			}

			if _, err = io.Copy(mimePart, file); err != nil {
				return nil, err
			}

			if err = file.Close(); err != nil {
				return nil, err
			}
		}
	}

	for key, values := range form.fields {
		for _, value := range values {
			if err := bodyWriter.WriteField(key, value); err != nil {
				return nil, err
			}
		}
	}

	if err := bodyWriter.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(verb, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header = form.Header()
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	return req, nil
}

func postWithoutFiles(form *Form, verb, urlStr string) (*http.Request, error) {
	data := strings.NewReader(form.fields.Encode())
	req, err := http.NewRequest(verb, urlStr, data)
	if err != nil {
		return nil, err
	}
	req.Header = form.Header()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func getRequest(form *Form, verb, urlStr string) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	urlStr = u.Scheme +
		"://" +
		u.Host +
		u.Path +
		"?" +
		u.Query().Encode() + "&" +
		form.fields.Encode()

	req, err := http.NewRequest(verb, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header = form.Header()

	return req, nil
}
