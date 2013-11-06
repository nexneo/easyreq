package easyreq

import (
	"bytes"
	"fmt"
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
}

func NewForm(fields, files url.Values) (f *Form) {
	if fields != nil {
		f.fields = fields
	}
	if files != nil {
		f.files = files
	}
	return
}

func (f *Form) Fields() url.Values {
	if f.fields == nil {
		f.fields = make(url.Values)
	}
	return f.fields
}

func (f *Form) Files() url.Values {
	if f.files == nil {
		f.files = make(url.Values)
	}
	return f.files
}

func (form *Form) Request(verb, urlStr string) (*http.Request, error) {
	if verb == "GET" {
		return nil, fmt.Errorf("Can't create GET form [TODO]: %s", urlStr)
	}

	if len(form.files) == 0 {
		data := strings.NewReader(form.fields.Encode())
		req, _ := http.NewRequest(verb, urlStr, data)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return req, nil
	}

	body := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(body)

	for key, paths := range form.files {
		for _, pathVal := range paths {
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
			file.Close()
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
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	return req, nil
}
