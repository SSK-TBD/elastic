// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"fmt"
	"net/http"
)

// OpenPointInTimeService opens a point in time that can be used in subsequent
// searches.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.x/point-in-time-api.html
// for details.
type OpenPointInTimeService struct {
	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	index             []string
	preference        string
	routing           string
	ignoreUnavailable *bool
	expandWildcards   string
	keepAlive         string
	bodyJson          interface{}
	bodyString        string
}

// NewOpenPointInTimeService creates a new OpenPointInTimeService.
func NewOpenPointInTimeService() *OpenPointInTimeService {
	return &OpenPointInTimeService{}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *OpenPointInTimeService) Pretty(pretty bool) *OpenPointInTimeService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *OpenPointInTimeService) Human(human bool) *OpenPointInTimeService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *OpenPointInTimeService) ErrorTrace(errorTrace bool) *OpenPointInTimeService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *OpenPointInTimeService) FilterPath(filterPath ...string) *OpenPointInTimeService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *OpenPointInTimeService) Header(name string, value string) *OpenPointInTimeService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *OpenPointInTimeService) Headers(headers http.Header) *OpenPointInTimeService {
	s.headers = headers
	return s
}

// Preference specifies the node or shard the operation should be performed on.
func (s *OpenPointInTimeService) Preference(preference string) *OpenPointInTimeService {
	s.preference = preference
	return s
}

// Index is the name of the index (or indices).
func (s *OpenPointInTimeService) Index(index ...string) *OpenPointInTimeService {
	s.index = index
	return s
}

// Routing is a specific routing value.
func (s *OpenPointInTimeService) Routing(routing string) *OpenPointInTimeService {
	s.routing = routing
	return s
}

// IgnoreUnavailable indicates whether specified concrete indices should be
// ignored when unavailable (missing or closed).
func (s *OpenPointInTimeService) IgnoreUnavailable(ignoreUnavailable bool) *OpenPointInTimeService {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// ExpandWildcards indicates whether to expand wildcard expression to
// concrete indices that are open, closed or both.
func (s *OpenPointInTimeService) ExpandWildcards(expandWildcards string) *OpenPointInTimeService {
	s.expandWildcards = expandWildcards
	return s
}

// KeepAlive indicates the specific time to live for the point in time.
func (s *OpenPointInTimeService) KeepAlive(keepAlive string) *OpenPointInTimeService {
	s.keepAlive = keepAlive
	return s
}

// BodyJson is the document as a serializable JSON interface.
func (s *OpenPointInTimeService) BodyJson(body interface{}) *OpenPointInTimeService {
	s.bodyJson = body
	return s
}

// BodyString is the document encoded as a string.
func (s *OpenPointInTimeService) BodyString(body string) *OpenPointInTimeService {
	s.bodyString = body
	return s
}

// Validate checks if the operation is valid.
func (s *OpenPointInTimeService) Validate() error {
	var invalid []string
	if len(s.index) == 0 {
		invalid = append(invalid, "Index")
	}
	if s.keepAlive == "" {
		invalid = append(invalid, "KeepAlive")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}
