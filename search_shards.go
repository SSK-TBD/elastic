// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/http"
	"time"
)

// SearchShardsService returns the indices and shards that a search request would be executed against.
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/search-shards.html
type SearchShardsService struct {
	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	index             []string
	routing           string
	local             *bool
	preference        string
	ignoreUnavailable *bool
	allowNoIndices    *bool
	expandWildcards   string
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *SearchShardsService) Pretty(pretty bool) *SearchShardsService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SearchShardsService) Human(human bool) *SearchShardsService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SearchShardsService) ErrorTrace(errorTrace bool) *SearchShardsService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SearchShardsService) FilterPath(filterPath ...string) *SearchShardsService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SearchShardsService) Header(name string, value string) *SearchShardsService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SearchShardsService) Headers(headers http.Header) *SearchShardsService {
	s.headers = headers
	return s
}

// Index sets the names of the indices to restrict the results.
func (s *SearchShardsService) Index(index ...string) *SearchShardsService {
	s.index = append(s.index, index...)
	return s
}

// A boolean value whether to read the cluster state locally in order to
// determine where shards are allocated instead of using the Master node’s cluster state.
func (s *SearchShardsService) Local(local bool) *SearchShardsService {
	s.local = &local
	return s
}

// Routing sets a specific routing value.
func (s *SearchShardsService) Routing(routing string) *SearchShardsService {
	s.routing = routing
	return s
}

// Preference specifies the node or shard the operation should be performed on (default: random).
func (s *SearchShardsService) Preference(preference string) *SearchShardsService {
	s.preference = preference
	return s
}

// IgnoreUnavailable indicates whether the specified concrete indices
// should be ignored when unavailable (missing or closed).
func (s *SearchShardsService) IgnoreUnavailable(ignoreUnavailable bool) *SearchShardsService {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// AllowNoIndices indicates whether to ignore if a wildcard indices
// expression resolves into no concrete indices. (This includes `_all` string
// or when no indices have been specified).
func (s *SearchShardsService) AllowNoIndices(allowNoIndices bool) *SearchShardsService {
	s.allowNoIndices = &allowNoIndices
	return s
}

// ExpandWildcards indicates whether to expand wildcard expression to
// concrete indices that are open, closed or both.
func (s *SearchShardsService) ExpandWildcards(expandWildcards string) *SearchShardsService {
	s.expandWildcards = expandWildcards
	return s
}

type RecoverySource struct {
	Type string `json:"type"`
	// TODO add missing fields here based on the Type
}

type AllocationId struct {
	Id           string `json:"id"`
	RelocationId string `json:"relocation_id,omitempty"`
}

type UnassignedInfo struct {
	Reason           string     `json:"reason"`
	At               *time.Time `json:"at,omitempty"`
	FailedAttempts   int        `json:"failed_attempts,omitempty"`
	Delayed          bool       `json:"delayed"`
	Details          string     `json:"details,omitempty"`
	AllocationStatus string     `json:"allocation_status"`
}
