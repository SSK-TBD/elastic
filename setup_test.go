// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	testIndexName      = "elastic-test"
	testIndexName2     = "elastic-test2"
	testIndexName3     = "elastic-test3"
	testIndexName4     = "elastic-test4"
	testIndexName5     = "elastic-test5"
	testIndexNameEmpty = "elastic-test-empty"
	testMapping        = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"tags":{
				"type":"keyword"
			},
			"location":{
				"type":"geo_point"
			},
			"suggest_field":{
				"type":"completion"
			}
		}
	}
}
`
	testMappingWithContext = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"tags":{
				"type":"keyword"
			},
			"location":{
				"type":"geo_point"
			},
			"suggest_field":{
				"type":"completion",
				"contexts":[
					{
						"name":"user_name",
						"type":"category"
					}
				]
			}
		}
	}
}
`

	testNoSourceIndexName = "elastic-nosource-test"
	testNoSourceMapping   = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"_source": {
			"enabled": false
		},
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"tags":{
				"type":"keyword"
			},
			"location":{
				"type":"geo_point"
			},
			"suggest_field":{
				"type":"completion",
				"contexts":[
					{
						"name":"user_name",
						"type":"category"
					}
				]
			}
		}
	}
}
`

	testJoinIndex   = "elastic-joins"
	testJoinMapping = `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"message":{
					"type":"text"
				},
				"my_join_field": {
					"type": "join",
					"relations": {
						"question": "answer"
					}
				}
			}
		}
	}
`

	testOrderIndex   = "elastic-orders"
	testOrderMapping = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"article":{
				"type":"text"
			},
			"manufacturer":{
				"type":"keyword"
			},
			"price":{
				"type":"float"
			},
			"time":{
				"type":"date",
				"format": "yyyy-MM-dd"
			}
		}
	}
}
`

	/*
		   	testDoctypeIndex   = "elastic-doctypes"
		   	testDoctypeMapping = `
		   {
		   	"settings":{
		   		"number_of_shards":1,
		   		"number_of_replicas":0
		   	},
		   	"mappings":{
				"properties":{
					"message":{
						"type":"text",
						"store": true,
						"fielddata": true
					}
				}
		   	}
		   }
		   `
	*/

	testQueryIndex   = "elastic-queries"
	testQueryMapping = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"message":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"query": {
				"type":	"percolator"
			}
		}
	}
}
`
)

type tweet struct {
	User     string        `json:"user"`
	Message  string        `json:"message"`
	Retweets int           `json:"retweets"`
	Image    string        `json:"image,omitempty"`
	Created  time.Time     `json:"created,omitempty"`
	Tags     []string      `json:"tags,omitempty"`
	Location string        `json:"location,omitempty"`
	Suggest  *SuggestField `json:"suggest_field,omitempty"`
}

func (t tweet) String() string {
	return fmt.Sprintf("tweet{User:%q,Message:%q,Retweets:%d}", t.User, t.Message, t.Retweets)
}

type joinDoc struct {
	Message   string      `json:"message"`
	JoinField interface{} `json:"my_join_field,omitempty"`
}

type joinField struct {
	Name   string `json:"name"`
	Parent string `json:"parent,omitempty"`
}

type order struct {
	Article      string  `json:"article"`
	Manufacturer string  `json:"manufacturer"`
	Price        float64 `json:"price"`
	Time         string  `json:"time,omitempty"`
}

func (o order) String() string {
	return fmt.Sprintf("order{Article:%q,Manufacturer:%q,Price:%v,Time:%v}", o.Article, o.Manufacturer, o.Price, o.Time)
}

// doctype is required for Percolate tests.
type doctype struct {
	Message string `json:"message"`
}

func isCI() bool {
	return os.Getenv("TRAVIS") != "" || os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
}

type logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fail()
	FailNow()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

func boolPtr(b bool) *bool { return &b }

// strictDecoder returns an error if any JSON fields aren't decoded.
type strictDecoder struct{}

func (d *strictDecoder) Decode(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

var (
	logDeprecations = flag.String("deprecations", "off", "log or fail on deprecation warnings")
	logTypesRemoval = flag.Bool("types-removal", false, "log deprecation warnings regarding types removal")
	strict          = flag.Bool("strict-decoder", false, "treat missing unknown fields in response as errors")
	noHealthcheck   = flag.Bool("no-healthcheck", false, "allows to disable healthchecks globally")
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type lexicographically struct {
	strings []string
}

func (l lexicographically) Len() int {
	return len(l.strings)
}

func (l lexicographically) Less(i, j int) bool {
	return l.strings[i] < l.strings[j]
}

func (l lexicographically) Swap(i, j int) {
	l.strings[i], l.strings[j] = l.strings[j], l.strings[i]
}
