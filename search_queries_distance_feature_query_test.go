// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestDistanceFeatureQueryForDateField(t *testing.T) {
	q := NewDistanceFeatureQuery("production_date", "now", "7d")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"distance_feature":{"field":"production_date","origin":"now","pivot":"7d"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestDistanceFeatureQueryForGeoField(t *testing.T) {
	q := NewDistanceFeatureQuery("location", GeoPointFromLatLon(-71.3, 41.15), "1000m")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"distance_feature":{"field":"location","origin":{"lat":-71.3,"lon":41.15},"pivot":"1000m"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
