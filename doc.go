// Package easyreq provides support for creating requests easily for
// multipart form request or json API requests.
//
// Usage
// import "github.com/nexneo/easyreq"
//
// Form Example
//		f := easyreq.Form{}
//		f.Field().Add("Name", "John")
//		f.File().Add("File", "test-files/logo.png")
//		f.Header().Add("Host", "example.com")
//
// 		req, err := f.Request("POST", "http://example.com/postform")
// or
//		easyreq.NewForm(fields, nil).Do("GET", "http://example.com/")
//		easyreq.NewForm(fields, files).Do("POST", "http://example.com/")
//
// Form will choose Content-Type based on any file added or not.
//
// Json Example
//		j := easyreq.Json{}
//		j.Header().Add("Host", "example.com")
// 		req, err := j.Set(v).Request("POST", "http://example.com/postjson")
// or
// 		req, err := easyreq.NewJson(v).Request("PUT", "http://example.com/putjson")
package easyreq
