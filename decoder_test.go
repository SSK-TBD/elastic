// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"encoding/json"
	"sync/atomic"
)

type decoder struct {
	N int64
}

func (d *decoder) Decode(data []byte, v interface{}) error {
	atomic.AddInt64(&d.N, 1)
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(v)
}
