// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/http"
)

// GetScriptService reads a stored script in Elasticsearch.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/modules-scripting.html
// for details.
type GetScriptService struct {
	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	id string
}

// NewGetScriptService creates a new GetScriptService.
func NewGetScriptService() *GetScriptService {
	return &GetScriptService{}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *GetScriptService) Pretty(pretty bool) *GetScriptService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *GetScriptService) Human(human bool) *GetScriptService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *GetScriptService) ErrorTrace(errorTrace bool) *GetScriptService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *GetScriptService) FilterPath(filterPath ...string) *GetScriptService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *GetScriptService) Header(name string, value string) *GetScriptService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *GetScriptService) Headers(headers http.Header) *GetScriptService {
	s.headers = headers
	return s
}

// Id is the script ID.
func (s *GetScriptService) Id(id string) *GetScriptService {
	s.id = id
	return s
}
