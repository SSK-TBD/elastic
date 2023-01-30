// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic_test

import (
	"encoding/json"
	"testing"

	"github.com/SSK-TBD/elastic/v7"
)

func TestWildcardQuery(t *testing.T) {
	q := elastic.NewWildcardQuery("user", "ki*y??")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"wildcard":{"user":{"value":"ki*y??"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestWildcardQueryWithBoost(t *testing.T) {
	q := elastic.NewWildcardQuery("user", "ki*y??").Boost(1.2)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"wildcard":{"user":{"boost":1.2,"value":"ki*y??"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestWildcardQueryWithCaseInsensitive(t *testing.T) {
	q := elastic.NewWildcardQuery("user", "ki*y??").CaseInsensitive(true)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"wildcard":{"user":{"case_insensitive":true,"value":"ki*y??"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
