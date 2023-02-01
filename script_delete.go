// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/http"
)

// DeleteScriptService removes a stored script in Elasticsearch.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/modules-scripting.html
// for details.
type DeleteScriptService struct {
	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	id            string
	timeout       string
	masterTimeout string
}

// NewDeleteScriptService creates a new DeleteScriptService.
func NewDeleteScriptService() *DeleteScriptService {
	return &DeleteScriptService{}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *DeleteScriptService) Pretty(pretty bool) *DeleteScriptService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *DeleteScriptService) Human(human bool) *DeleteScriptService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *DeleteScriptService) ErrorTrace(errorTrace bool) *DeleteScriptService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *DeleteScriptService) FilterPath(filterPath ...string) *DeleteScriptService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *DeleteScriptService) Header(name string, value string) *DeleteScriptService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *DeleteScriptService) Headers(headers http.Header) *DeleteScriptService {
	s.headers = headers
	return s
}

// Id is the script ID.
func (s *DeleteScriptService) Id(id string) *DeleteScriptService {
	s.id = id
	return s
}

// Timeout is an explicit operation timeout.
func (s *DeleteScriptService) Timeout(timeout string) *DeleteScriptService {
	s.timeout = timeout
	return s
}

// MasterTimeout is the timeout for connecting to master.
func (s *DeleteScriptService) MasterTimeout(masterTimeout string) *DeleteScriptService {
	s.masterTimeout = masterTimeout
	return s
}
