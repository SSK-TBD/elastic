// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/http"
)

// PutScriptService adds or updates a stored script in Elasticsearch.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/modules-scripting.html
// for details.
type PutScriptService struct {
	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	id            string
	context       string
	timeout       string
	masterTimeout string
	bodyJson      interface{}
	bodyString    string
}

// NewPutScriptService creates a new PutScriptService.
func NewPutScriptService() *PutScriptService {
	return &PutScriptService{}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *PutScriptService) Pretty(pretty bool) *PutScriptService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *PutScriptService) Human(human bool) *PutScriptService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *PutScriptService) ErrorTrace(errorTrace bool) *PutScriptService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *PutScriptService) FilterPath(filterPath ...string) *PutScriptService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *PutScriptService) Header(name string, value string) *PutScriptService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *PutScriptService) Headers(headers http.Header) *PutScriptService {
	s.headers = headers
	return s
}

// Id is the script ID.
func (s *PutScriptService) Id(id string) *PutScriptService {
	s.id = id
	return s
}

// Context specifies the script context (optional).
func (s *PutScriptService) Context(context string) *PutScriptService {
	s.context = context
	return s
}

// Timeout is an explicit operation timeout.
func (s *PutScriptService) Timeout(timeout string) *PutScriptService {
	s.timeout = timeout
	return s
}

// MasterTimeout is the timeout for connecting to master.
func (s *PutScriptService) MasterTimeout(masterTimeout string) *PutScriptService {
	s.masterTimeout = masterTimeout
	return s
}

// BodyJson is the document as a serializable JSON interface.
func (s *PutScriptService) BodyJson(body interface{}) *PutScriptService {
	s.bodyJson = body
	return s
}

// BodyString is the document encoded as a string.
func (s *PutScriptService) BodyString(body string) *PutScriptService {
	s.bodyString = body
	return s
}
