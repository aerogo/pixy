package pixy

import (
	"bytes"
	"sync"
)

var pool sync.Pool

func acquireBytesBuffer() *bytes.Buffer {
	var _b *bytes.Buffer
	obj := pool.Get()

	if obj == nil {
		return &bytes.Buffer{}
	}

	_b = obj.(*bytes.Buffer)
	_b.Reset()
	return _b
}
