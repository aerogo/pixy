package pixy

import (
	"strings"
	"sync"
)

var pool sync.Pool

func acquireStringsBuilder() *strings.Builder {
	var _b *strings.Builder
	obj := pool.Get()

	if obj == nil {
		return &strings.Builder{}
	}

	_b = obj.(*strings.Builder)
	_b.Reset()
	return _b
}
