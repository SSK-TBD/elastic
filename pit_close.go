// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ClosePointInTimeService removes a point in time.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.x/point-in-time-api.html
// for details.
type ClosePointInTimeService struct {
	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	id         string
	bodyJson   interface{}
	bodyString string
}

// NewClosePointInTimeService creates a new ClosePointInTimeService.
func NewClosePointInTimeService() *ClosePointInTimeService {
	return &ClosePointInTimeService{}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *ClosePointInTimeService) Pretty(pretty bool) *ClosePointInTimeService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *ClosePointInTimeService) Human(human bool) *ClosePointInTimeService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *ClosePointInTimeService) ErrorTrace(errorTrace bool) *ClosePointInTimeService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *ClosePointInTimeService) FilterPath(filterPath ...string) *ClosePointInTimeService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *ClosePointInTimeService) Header(name string, value string) *ClosePointInTimeService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *ClosePointInTimeService) Headers(headers http.Header) *ClosePointInTimeService {
	s.headers = headers
	return s
}

// ID to close.
func (s *ClosePointInTimeService) ID(id string) *ClosePointInTimeService {
	s.id = id
	return s
}

// BodyJson is the document as a serializable JSON interface.
func (s *ClosePointInTimeService) BodyJson(body interface{}) *ClosePointInTimeService {
	s.bodyJson = body
	return s
}

// BodyString is the document encoded as a string.
func (s *ClosePointInTimeService) BodyString(body string) *ClosePointInTimeService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *ClosePointInTimeService) buildURL() (string, string, url.Values, error) {
	var (
		method = "DELETE"
		path   = "/_pit"
	)

	// Add query string parameters
	params := url.Values{}
	if v := s.pretty; v != nil {
		params.Set("pretty", fmt.Sprint(*v))
	}
	if v := s.human; v != nil {
		params.Set("human", fmt.Sprint(*v))
	}
	if v := s.errorTrace; v != nil {
		params.Set("error_trace", fmt.Sprint(*v))
	}
	if len(s.filterPath) > 0 {
		params.Set("filter_path", strings.Join(s.filterPath, ","))
	}
	return method, path, params, nil
}

// Validate checks if the operation is valid.
func (s *ClosePointInTimeService) Validate() error {
	return nil
}
