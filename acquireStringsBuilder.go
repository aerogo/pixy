package pixy

import (
	"strings"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

func acquireStringsBuilder() *strings.Builder {
	builder := pool.Get().(*strings.Builder)
	builder.Reset()
	return builder
}
