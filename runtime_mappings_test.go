// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestRuntimeMappingsSource(t *testing.T) {
	rm := RuntimeMappings{
		"day_of_week": map[string]interface{}{
			"type": "keyword",
		},
	}
	src, err := rm.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"day_of_week":{"type":"keyword"}}`
	if want, have := expected, string(data); want != have {
		t.Fatalf("want %s, have %s", want, have)
	}
}
