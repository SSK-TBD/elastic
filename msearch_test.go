// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// import (
// 	"context"
// 	"encoding/json"
// 	_ "net/http"
// 	"testing"
// )

// func TestMultiSearch(t *testing.T) {
// 	client := setupTestClientAndCreateIndex(t)
// 	// client := setupTestClientAndCreateIndexAndLog(t)

// 	tweet1 := tweet{
// 		User:    "olivere",
// 		Message: "Welcome to Golang and Elasticsearch.",
// 		Tags:    []string{"golang", "elasticsearch"},
// 	}
// 	tweet2 := tweet{
// 		User:    "olivere",
// 		Message: "Another unrelated topic.",
// 		Tags:    []string{"golang"},
// 	}
// 	tweet3 := tweet{
// 		User:    "sandrae",
// 		Message: "Cycling is fun.",
// 		Tags:    []string{"sports", "cycling"},
// 	}

// 	// Add all documents
// 	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Index().Index(testIndexName).Id("3").BodyJson(&tweet3).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Refresh().Index(testIndexName).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Spawn two search queries with one roundtrip
// 	q1 := NewMatchAllQuery()
// 	q2 := NewTermQuery("tags", "golang")

// 	sreq1 := NewSearchRequest().Index(testIndexName, testIndexName2).
// 		Source(NewSearchSource().Query(q1).Size(10))
// 	sreq2 := NewSearchRequest().Index(testIndexName).
// 		Source(NewSearchSource().Query(q2))

// 	searchResult, err := client.MultiSearch().
// 		Add(sreq1, sreq2).
// 		Pretty(true).
// 		Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if searchResult.Responses == nil {
// 		t.Fatal("expected responses != nil; got nil")
// 	}
// 	if len(searchResult.Responses) != 2 {
// 		t.Fatalf("expected 2 responses; got %d", len(searchResult.Responses))
// 	}

// 	sres := searchResult.Responses[0]
// 	if sres.Hits == nil {
// 		t.Errorf("expected Hits != nil; got nil")
// 	}
// 	if sres.TotalHits() != 3 {
// 		t.Errorf("expected TotalHits() = %d; got %d", 3, sres.TotalHits())
// 	}
// 	if len(sres.Hits.Hits) != 3 {
// 		t.Errorf("expected len(Hits.Hits) = %d; got %d", 3, len(sres.Hits.Hits))
// 	}
// 	for _, hit := range sres.Hits.Hits {
// 		if hit.Index != testIndexName {
// 			t.Errorf("expected Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
// 		}
// 		item := make(map[string]interface{})
// 		err := json.Unmarshal(hit.Source, &item)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	sres = searchResult.Responses[1]
// 	if sres.Hits == nil {
// 		t.Errorf("expected Hits != nil; got nil")
// 	}
// 	if sres.TotalHits() != 2 {
// 		t.Errorf("expected TotalHits() = %d; got %d", 2, sres.TotalHits())
// 	}
// 	if len(sres.Hits.Hits) != 2 {
// 		t.Errorf("expected len(Hits.Hits) = %d; got %d", 2, len(sres.Hits.Hits))
// 	}
// 	for _, hit := range sres.Hits.Hits {
// 		if hit.Index != testIndexName {
// 			t.Errorf("expected Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
// 		}
// 		item := make(map[string]interface{})
// 		err := json.Unmarshal(hit.Source, &item)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// }

// func TestMultiSearchWithStrings(t *testing.T) {
// 	client := setupTestClientAndCreateIndex(t)
// 	// client := setupTestClientAndCreateIndexAndLog(t)

// 	tweet1 := tweet{
// 		User:    "olivere",
// 		Message: "Welcome to Golang and Elasticsearch.",
// 		Tags:    []string{"golang", "elasticsearch"},
// 	}
// 	tweet2 := tweet{
// 		User:    "olivere",
// 		Message: "Another unrelated topic.",
// 		Tags:    []string{"golang"},
// 	}
// 	tweet3 := tweet{
// 		User:    "sandrae",
// 		Message: "Cycling is fun.",
// 		Tags:    []string{"sports", "cycling"},
// 	}

// 	// Add all documents
// 	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Index().Index(testIndexName).Id("3").BodyJson(&tweet3).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Refresh().Index(testIndexName).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Spawn two search queries with one roundtrip
// 	sreq1 := NewSearchRequest().Index(testIndexName, testIndexName2).
// 		Source(`{"query":{"match_all":{}}}`)
// 	sreq2 := NewSearchRequest().Index(testIndexName).
// 		Source(`{"query":{"term":{"tags":"golang"}}}`)

// 	searchResult, err := client.MultiSearch().
// 		Add(sreq1, sreq2).
// 		Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if searchResult.Responses == nil {
// 		t.Fatal("expected responses != nil; got nil")
// 	}
// 	if len(searchResult.Responses) != 2 {
// 		t.Fatalf("expected 2 responses; got %d", len(searchResult.Responses))
// 	}

// 	sres := searchResult.Responses[0]
// 	if sres.Hits == nil {
// 		t.Errorf("expected Hits != nil; got nil")
// 	}
// 	if sres.TotalHits() != 3 {
// 		t.Errorf("expected TotalHits() = %d; got %d", 3, sres.TotalHits())
// 	}
// 	if len(sres.Hits.Hits) != 3 {
// 		t.Errorf("expected len(Hits.Hits) = %d; got %d", 3, len(sres.Hits.Hits))
// 	}
// 	for _, hit := range sres.Hits.Hits {
// 		if hit.Index != testIndexName {
// 			t.Errorf("expected Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
// 		}
// 		item := make(map[string]interface{})
// 		err := json.Unmarshal(hit.Source, &item)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	sres = searchResult.Responses[1]
// 	if sres.Hits == nil {
// 		t.Errorf("expected Hits != nil; got nil")
// 	}
// 	if sres.TotalHits() != 2 {
// 		t.Errorf("expected TotalHits() = %d; got %d", 2, sres.TotalHits())
// 	}
// 	if len(sres.Hits.Hits) != 2 {
// 		t.Errorf("expected len(Hits.Hits) = %d; got %d", 2, len(sres.Hits.Hits))
// 	}
// 	for _, hit := range sres.Hits.Hits {
// 		if hit.Index != testIndexName {
// 			t.Errorf("expected Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
// 		}
// 		item := make(map[string]interface{})
// 		err := json.Unmarshal(hit.Source, &item)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// }

// func TestMultiSearchWithOneRequest(t *testing.T) {
// 	client := setupTestClientAndCreateIndex(t)

// 	tweet1 := tweet{
// 		User:    "olivere",
// 		Message: "Welcome to Golang and Elasticsearch.",
// 		Tags:    []string{"golang", "elasticsearch"},
// 	}
// 	tweet2 := tweet{
// 		User:    "olivere",
// 		Message: "Another unrelated topic.",
// 		Tags:    []string{"golang"},
// 	}
// 	tweet3 := tweet{
// 		User:    "sandrae",
// 		Message: "Cycling is fun.",
// 		Tags:    []string{"sports", "cycling"},
// 	}

// 	// Add all documents
// 	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Index().Index(testIndexName).Id("3").BodyJson(&tweet3).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Refresh().Index(testIndexName).Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Spawn two search queries with one roundtrip
// 	query := NewMatchAllQuery()
// 	source := NewSearchSource().Query(query).Size(10)
// 	sreq := NewSearchRequest().Source(source)

// 	searchResult, err := client.MultiSearch().
// 		Index(testIndexName).
// 		Add(sreq).
// 		Do(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if searchResult.Responses == nil {
// 		t.Fatal("expected responses != nil; got nil")
// 	}
// 	if len(searchResult.Responses) != 1 {
// 		t.Fatalf("expected 1 responses; got %d", len(searchResult.Responses))
// 	}

// 	sres := searchResult.Responses[0]
// 	if sres.Hits == nil {
// 		t.Errorf("expected Hits != nil; got nil")
// 	}
// 	if sres.TotalHits() != 3 {
// 		t.Errorf("expected TotalHits() = %d; got %d", 3, sres.TotalHits())
// 	}
// 	if len(sres.Hits.Hits) != 3 {
// 		t.Errorf("expected len(Hits.Hits) = %d; got %d", 3, len(sres.Hits.Hits))
// 	}
// 	for _, hit := range sres.Hits.Hits {
// 		if hit.Index != testIndexName {
// 			t.Errorf("expected Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
// 		}
// 		item := make(map[string]interface{})
// 		err := json.Unmarshal(hit.Source, &item)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// }
