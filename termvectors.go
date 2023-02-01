// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/http"
)

// TermvectorsService returns information and statistics on terms in the
// fields of a particular document. The document could be stored in the
// index or artificially provided by the user.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/docs-termvectors.html
// for documentation.
type TermvectorsService struct {
	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	id               string
	index            string
	typ              string
	dfs              *bool
	doc              interface{}
	fieldStatistics  *bool
	fields           []string
	filter           *TermvectorsFilterSettings
	perFieldAnalyzer map[string]string
	offsets          *bool
	parent           string
	payloads         *bool
	positions        *bool
	preference       string
	realtime         *bool
	routing          string
	termStatistics   *bool
	version          interface{}
	versionType      string
	bodyJson         interface{}
	bodyString       string
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *TermvectorsService) Pretty(pretty bool) *TermvectorsService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *TermvectorsService) Human(human bool) *TermvectorsService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *TermvectorsService) ErrorTrace(errorTrace bool) *TermvectorsService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *TermvectorsService) FilterPath(filterPath ...string) *TermvectorsService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *TermvectorsService) Header(name string, value string) *TermvectorsService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *TermvectorsService) Headers(headers http.Header) *TermvectorsService {
	s.headers = headers
	return s
}

// Index in which the document resides.
func (s *TermvectorsService) Index(index string) *TermvectorsService {
	s.index = index
	return s
}

// Type of the document.
//
// Deprecated: Types are in the process of being removed.
func (s *TermvectorsService) Type(typ string) *TermvectorsService {
	s.typ = typ
	return s
}

// Id of the document.
func (s *TermvectorsService) Id(id string) *TermvectorsService {
	s.id = id
	return s
}

// Dfs specifies if distributed frequencies should be returned instead
// shard frequencies.
func (s *TermvectorsService) Dfs(dfs bool) *TermvectorsService {
	s.dfs = &dfs
	return s
}

// Doc is the document to analyze.
func (s *TermvectorsService) Doc(doc interface{}) *TermvectorsService {
	s.doc = doc
	return s
}

// FieldStatistics specifies if document count, sum of document frequencies
// and sum of total term frequencies should be returned.
func (s *TermvectorsService) FieldStatistics(fieldStatistics bool) *TermvectorsService {
	s.fieldStatistics = &fieldStatistics
	return s
}

// Fields a list of fields to return.
func (s *TermvectorsService) Fields(fields ...string) *TermvectorsService {
	if s.fields == nil {
		s.fields = make([]string, 0)
	}
	s.fields = append(s.fields, fields...)
	return s
}

// Filter adds terms filter settings.
func (s *TermvectorsService) Filter(filter *TermvectorsFilterSettings) *TermvectorsService {
	s.filter = filter
	return s
}

// PerFieldAnalyzer allows to specify a different analyzer than the one
// at the field.
func (s *TermvectorsService) PerFieldAnalyzer(perFieldAnalyzer map[string]string) *TermvectorsService {
	s.perFieldAnalyzer = perFieldAnalyzer
	return s
}

// Offsets specifies if term offsets should be returned.
func (s *TermvectorsService) Offsets(offsets bool) *TermvectorsService {
	s.offsets = &offsets
	return s
}

// Parent id of documents.
func (s *TermvectorsService) Parent(parent string) *TermvectorsService {
	s.parent = parent
	return s
}

// Payloads specifies if term payloads should be returned.
func (s *TermvectorsService) Payloads(payloads bool) *TermvectorsService {
	s.payloads = &payloads
	return s
}

// Positions specifies if term positions should be returned.
func (s *TermvectorsService) Positions(positions bool) *TermvectorsService {
	s.positions = &positions
	return s
}

// Preference specify the node or shard the operation
// should be performed on (default: random).
func (s *TermvectorsService) Preference(preference string) *TermvectorsService {
	s.preference = preference
	return s
}

// Realtime specifies if request is real-time as opposed to
// near-real-time (default: true).
func (s *TermvectorsService) Realtime(realtime bool) *TermvectorsService {
	s.realtime = &realtime
	return s
}

// Routing is a specific routing value.
func (s *TermvectorsService) Routing(routing string) *TermvectorsService {
	s.routing = routing
	return s
}

// TermStatistics specifies if total term frequency and document frequency
// should be returned.
func (s *TermvectorsService) TermStatistics(termStatistics bool) *TermvectorsService {
	s.termStatistics = &termStatistics
	return s
}

// Version an explicit version number for concurrency control.
func (s *TermvectorsService) Version(version interface{}) *TermvectorsService {
	s.version = version
	return s
}

// VersionType specifies a version type ("internal", "external", or "external_gte").
func (s *TermvectorsService) VersionType(versionType string) *TermvectorsService {
	s.versionType = versionType
	return s
}

// BodyJson defines the body parameters. See documentation.
func (s *TermvectorsService) BodyJson(body interface{}) *TermvectorsService {
	s.bodyJson = body
	return s
}

// BodyString defines the body parameters as a string. See documentation.
func (s *TermvectorsService) BodyString(body string) *TermvectorsService {
	s.bodyString = body
	return s
}

// -- Filter settings --

// TermvectorsFilterSettings adds additional filters to a Termsvector request.
// It allows to filter terms based on their tf-idf scores.
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/docs-termvectors.html#_terms_filtering
// for more information.
type TermvectorsFilterSettings struct {
	maxNumTerms   *int64
	minTermFreq   *int64
	maxTermFreq   *int64
	minDocFreq    *int64
	maxDocFreq    *int64
	minWordLength *int64
	maxWordLength *int64
}

// NewTermvectorsFilterSettings creates and initializes a new TermvectorsFilterSettings struct.
func NewTermvectorsFilterSettings() *TermvectorsFilterSettings {
	return &TermvectorsFilterSettings{}
}

// MaxNumTerms specifies the maximum number of terms the must be returned per field.
func (fs *TermvectorsFilterSettings) MaxNumTerms(value int64) *TermvectorsFilterSettings {
	fs.maxNumTerms = &value
	return fs
}

// MinTermFreq ignores words with less than this frequency in the source doc.
func (fs *TermvectorsFilterSettings) MinTermFreq(value int64) *TermvectorsFilterSettings {
	fs.minTermFreq = &value
	return fs
}

// MaxTermFreq ignores words with more than this frequency in the source doc.
func (fs *TermvectorsFilterSettings) MaxTermFreq(value int64) *TermvectorsFilterSettings {
	fs.maxTermFreq = &value
	return fs
}

// MinDocFreq ignores terms which do not occur in at least this many docs.
func (fs *TermvectorsFilterSettings) MinDocFreq(value int64) *TermvectorsFilterSettings {
	fs.minDocFreq = &value
	return fs
}

// MaxDocFreq ignores terms which occur in more than this many docs.
func (fs *TermvectorsFilterSettings) MaxDocFreq(value int64) *TermvectorsFilterSettings {
	fs.maxDocFreq = &value
	return fs
}

// MinWordLength specifies the minimum word length below which words will be ignored.
func (fs *TermvectorsFilterSettings) MinWordLength(value int64) *TermvectorsFilterSettings {
	fs.minWordLength = &value
	return fs
}

// MaxWordLength specifies the maximum word length above which words will be ignored.
func (fs *TermvectorsFilterSettings) MaxWordLength(value int64) *TermvectorsFilterSettings {
	fs.maxWordLength = &value
	return fs
}

// Source returns JSON for the query.
func (fs *TermvectorsFilterSettings) Source() (interface{}, error) {
	source := make(map[string]interface{})
	if fs.maxNumTerms != nil {
		source["max_num_terms"] = *fs.maxNumTerms
	}
	if fs.minTermFreq != nil {
		source["min_term_freq"] = *fs.minTermFreq
	}
	if fs.maxTermFreq != nil {
		source["max_term_freq"] = *fs.maxTermFreq
	}
	if fs.minDocFreq != nil {
		source["min_doc_freq"] = *fs.minDocFreq
	}
	if fs.maxDocFreq != nil {
		source["max_doc_freq"] = *fs.maxDocFreq
	}
	if fs.minWordLength != nil {
		source["min_word_length"] = *fs.minWordLength
	}
	if fs.maxWordLength != nil {
		source["max_word_length"] = *fs.maxWordLength
	}
	return source, nil
}

// -- Response types --

type TokenInfo struct {
	StartOffset int64  `json:"start_offset"`
	EndOffset   int64  `json:"end_offset"`
	Position    int64  `json:"position"`
	Payload     string `json:"payload"`
}

type TermsInfo struct {
	DocFreq  int64       `json:"doc_freq"`
	Score    float64     `json:"score"`
	TermFreq int64       `json:"term_freq"`
	Ttf      int64       `json:"ttf"`
	Tokens   []TokenInfo `json:"tokens"`
}

type FieldStatistics struct {
	DocCount   int64 `json:"doc_count"`
	SumDocFreq int64 `json:"sum_doc_freq"`
	SumTtf     int64 `json:"sum_ttf"`
}

type TermVectorsFieldInfo struct {
	FieldStatistics FieldStatistics      `json:"field_statistics"`
	Terms           map[string]TermsInfo `json:"terms"`
}

// TermvectorsResponse is the response of TermvectorsService.Do.
type TermvectorsResponse struct {
	Index       string                          `json:"_index"`
	Type        string                          `json:"_type"`
	Id          string                          `json:"_id,omitempty"`
	Version     int                             `json:"_version"`
	Found       bool                            `json:"found"`
	Took        int64                           `json:"took"`
	TermVectors map[string]TermVectorsFieldInfo `json:"term_vectors"`
}
