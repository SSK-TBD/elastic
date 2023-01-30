// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestScriptScoreQuery(t *testing.T) {
	q := NewScriptScoreQuery(
		NewMatchQuery("message", "elasticsearch"),
		NewScript("doc['likes'].value / 10"),
	).MinScore(1.1).Boost(5.0).QueryName("my_query")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"script_score":{"_name":"my_query","boost":5,"min_score":1.1,"query":{"match":{"message":{"query":"elasticsearch"}}},"script":{"source":"doc['likes'].value / 10"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
