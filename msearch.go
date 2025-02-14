// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/http"
)

// MultiSearch executes one or more searches in one roundtrip.
type MultiSearchService struct {
	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	requests              []*SearchRequest
	indices               []string
	maxConcurrentRequests *int
	preFilterShardSize    *int
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *MultiSearchService) Pretty(pretty bool) *MultiSearchService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *MultiSearchService) Human(human bool) *MultiSearchService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *MultiSearchService) ErrorTrace(errorTrace bool) *MultiSearchService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *MultiSearchService) FilterPath(filterPath ...string) *MultiSearchService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *MultiSearchService) Header(name string, value string) *MultiSearchService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *MultiSearchService) Headers(headers http.Header) *MultiSearchService {
	s.headers = headers
	return s
}

func (s *MultiSearchService) Add(requests ...*SearchRequest) *MultiSearchService {
	s.requests = append(s.requests, requests...)
	return s
}

func (s *MultiSearchService) Index(indices ...string) *MultiSearchService {
	s.indices = append(s.indices, indices...)
	return s
}

func (s *MultiSearchService) MaxConcurrentSearches(max int) *MultiSearchService {
	s.maxConcurrentRequests = &max
	return s
}

func (s *MultiSearchService) PreFilterShardSize(size int) *MultiSearchService {
	s.preFilterShardSize = &size
	return s
}
