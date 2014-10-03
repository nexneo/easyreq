package easyreq

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Structure that contains form fields and upload files,
// can optionally have request headers.
type Form struct {
	fields url.Values
	files  url.Values
	header http.Header
}

// Inteface implemeted by Form, Json structures
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

// Helper to set Basic Auth header for request
func (f *Form) SetBasicAuth(username, password string) {
	setBasicAuth(f, username, password)
}

func setBasicAuth(r Requester, username, password string) {
	s := username + ":" + password
	r.Header().Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s)))
}

// Helper funcation send requests using http.DefaultClient
func (f *Form) Do(verb, urlStr string) (*http.Response, error) {
	return do(f, verb, urlStr)
}

func do(r Requester, verb, urlStr string) (*http.Response, error) {
	req, err := r.Request(verb, urlStr)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// Returns request based on current Fields and Files assoicated with form
// Request will always have correct Content-Type set for POST and PUT
func (form *Form) Request(verb, urlStr string) (*http.Request, error) {
	// GET handled differently then POST, PUT
	if verb == "GET" {
		return getRequest(form, verb, urlStr)
	}

	// no files take shortcut
	if len(form.files) == 0 {
		return postWithoutFiles(form, verb, urlStr)
	}

	body, err := ioutil.TempFile(os.TempDir(), "easyreq-request")
	if err != nil {
		return nil, err
	}
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
			if err = bodyWriter.WriteField(key, value); err != nil {
				return nil, err
			}
		}
	}

	if err = bodyWriter.Close(); err != nil {
		return nil, err
	}
	if err = body.Close(); err != nil {
		return nil, err
	}
	log.Println(body.Name())
	body, err = os.Open(body.Name())
	req, err := http.NewRequest(verb, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header = form.Header()
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	return req, nil
}

// send urlencoded fields in body
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

// send urlencoded fields in url
func getRequest(form *Form, verb, urlStr string) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// bother only if new fields are added
	if len(form.fields) != 0 {
		if len(u.Query()) != 0 {
			// if original url had encoded fields join
			urlStr += "&"
		} else {
			// else start
			urlStr += "?"
		}
		urlStr += form.fields.Encode()
	}

	req, err := http.NewRequest(verb, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header = form.Header()

	return req, nil
}
