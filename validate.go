// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/http"
)

// ValidateService allows a user to validate a potentially
// expensive query without executing it.
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/search-validate.html.
type ValidateService struct {
	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	index             []string
	typ               []string
	q                 string
	explain           *bool
	rewrite           *bool
	allShards         *bool
	lenient           *bool
	analyzer          string
	df                string
	analyzeWildcard   *bool
	defaultOperator   string
	ignoreUnavailable *bool
	allowNoIndices    *bool
	expandWildcards   string
	bodyJson          interface{}
	bodyString        string
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *ValidateService) Pretty(pretty bool) *ValidateService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *ValidateService) Human(human bool) *ValidateService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *ValidateService) ErrorTrace(errorTrace bool) *ValidateService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *ValidateService) FilterPath(filterPath ...string) *ValidateService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *ValidateService) Header(name string, value string) *ValidateService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *ValidateService) Headers(headers http.Header) *ValidateService {
	s.headers = headers
	return s
}

// Index sets the names of the indices to use for search.
func (s *ValidateService) Index(index ...string) *ValidateService {
	s.index = append(s.index, index...)
	return s
}

// Type adds search restrictions for a list of types.
//
// Deprecated: Types are in the process of being removed. Instead of using a type, prefer to
// filter on a field on the document.
func (s *ValidateService) Type(typ ...string) *ValidateService {
	s.typ = append(s.typ, typ...)
	return s
}

// Lenient specifies whether format-based query failures
// (such as providing text to a numeric field) should be ignored.
func (s *ValidateService) Lenient(lenient bool) *ValidateService {
	s.lenient = &lenient
	return s
}

// Query in the Lucene query string syntax.
func (s *ValidateService) Q(q string) *ValidateService {
	s.q = q
	return s
}

// An explain parameter can be specified to get more detailed information about why a query failed.
func (s *ValidateService) Explain(explain *bool) *ValidateService {
	s.explain = explain
	return s
}

// Provide a more detailed explanation showing the actual Lucene query that will be executed.
func (s *ValidateService) Rewrite(rewrite *bool) *ValidateService {
	s.rewrite = rewrite
	return s
}

// Execute validation on all shards instead of one random shard per index.
func (s *ValidateService) AllShards(allShards *bool) *ValidateService {
	s.allShards = allShards
	return s
}

// AnalyzeWildcard specifies whether wildcards and prefix queries
// in the query string query should be analyzed (default: false).
func (s *ValidateService) AnalyzeWildcard(analyzeWildcard bool) *ValidateService {
	s.analyzeWildcard = &analyzeWildcard
	return s
}

// Analyzer is the analyzer for the query string query.
func (s *ValidateService) Analyzer(analyzer string) *ValidateService {
	s.analyzer = analyzer
	return s
}

// Df is the default field for query string query (default: _all).
func (s *ValidateService) Df(df string) *ValidateService {
	s.df = df
	return s
}

// DefaultOperator is the default operator for query string query (AND or OR).
func (s *ValidateService) DefaultOperator(defaultOperator string) *ValidateService {
	s.defaultOperator = defaultOperator
	return s
}

// Query sets a query definition using the Query DSL.
func (s *ValidateService) Query(query Query) *ValidateService {
	src, err := query.Source()
	if err != nil {
		// Do nothing in case of an error
		return s
	}
	body := make(map[string]interface{})
	body["query"] = src
	s.bodyJson = body
	return s
}

// IgnoreUnavailable indicates whether the specified concrete indices
// should be ignored when unavailable (missing or closed).
func (s *ValidateService) IgnoreUnavailable(ignoreUnavailable bool) *ValidateService {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// AllowNoIndices indicates whether to ignore if a wildcard indices
// expression resolves into no concrete indices. (This includes `_all` string
// or when no indices have been specified).
func (s *ValidateService) AllowNoIndices(allowNoIndices bool) *ValidateService {
	s.allowNoIndices = &allowNoIndices
	return s
}

// ExpandWildcards indicates whether to expand wildcard expression to
// concrete indices that are open, closed or both.
func (s *ValidateService) ExpandWildcards(expandWildcards string) *ValidateService {
	s.expandWildcards = expandWildcards
	return s
}

// BodyJson sets the query definition using the Query DSL.
func (s *ValidateService) BodyJson(body interface{}) *ValidateService {
	s.bodyJson = body
	return s
}

// BodyString sets the query definition using the Query DSL as a string.
func (s *ValidateService) BodyString(body string) *ValidateService {
	s.bodyString = body
	return s
}
