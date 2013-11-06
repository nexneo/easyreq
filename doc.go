// Package easyreq provides support for creating requests easily for
// multipart form request or json API requests.
//
// Form Example
// 		f := Form{}
// 		f.Field().Add("Name", "John")
// 		f.File().Add("File", "test-files/logo.png")
//
// 		req, err := f.Request("POST", "http://example.com/postform")
//
// 		// It will choose Content-Type based on any file added or not.
//
// Json Example
// 		req, err := (&Json{}).Set(v).Request("POST", "http://example.com/postjson")
// or
// 		req, err := NewJson(v).Request("PUT", "http://example.com/putjson")
package easyreq
