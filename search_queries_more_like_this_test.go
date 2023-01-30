// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestMoreLikeThisQuerySourceWithLikeText(t *testing.T) {
	q := NewMoreLikeThisQuery().LikeText("Golang topic").Field("message")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	expected := `{"more_like_this":{"fields":["message"],"like":["Golang topic"]}}`
	if got != expected {
		t.Fatalf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMoreLikeThisQuerySourceWithLikeAndUnlikeItems(t *testing.T) {
	q := NewMoreLikeThisQuery()
	q = q.LikeItems(
		NewMoreLikeThisQueryItem().Id("1"),
		NewMoreLikeThisQueryItem().Index(testIndexName2).Type("comment").Id("2").Routing("routing_id"),
	)
	q = q.IgnoreLikeItems(NewMoreLikeThisQueryItem().Id("3"))
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	expected := `{"more_like_this":{"like":[{"_id":"1"},{"_id":"2","_index":"elastic-test2","_type":"comment","routing":"routing_id"}],"unlike":[{"_id":"3"}]}}`
	if got != expected {
		t.Fatalf("expected\n%s\n,got:\n%s", expected, got)
	}
}
