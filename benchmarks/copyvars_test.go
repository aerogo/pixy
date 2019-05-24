package pixy

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

const (
	_1 = "unknown string"
	_2 = "another string"
	_3 = "yet another string"
)

var (
	poolBuffers  sync.Pool
	poolBuilders sync.Pool
)

func acquireBytesBuffer() *bytes.Buffer {
	var _b *bytes.Buffer
	obj := poolBuffers.Get()

	if obj == nil {
		return &bytes.Buffer{}
	}

	_b = obj.(*bytes.Buffer)
	_b.Reset()
	return _b
}

func acquireStringsBuilder() *strings.Builder {
	var _b *strings.Builder
	obj := poolBuilders.Get()

	if obj == nil {
		return &strings.Builder{}
	}

	_b = obj.(*strings.Builder)
	_b.Reset()
	return _b
}

func BenchmarkA1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_b := make([]byte, len(_1)+len(_2)+len(_3))
		_l := 0

		_l += copy(_b[_l:], _1)
		_l += copy(_b[_l:], _2)
		_l += copy(_b[_l:], _3)

		_ = string(_b)
	}
}

func BenchmarkA2(b *testing.B) {
	_s := [3]string{}
	_i := 0

	_s[_i] = _1
	_i++
	_s[_i] = _2
	_i++
	_s[_i] = _3

	for n := 0; n < b.N; n++ {
		_b := make([]byte, len(_s[0])+len(_s[1])+len(_s[2]))
		_l := 0

		for i := 0; i < 3; i++ {
			_l += copy(_b[_l:], _s[i])
		}

		_ = string(_b)
	}
}

func BenchmarkA3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var buffer bytes.Buffer
		buffer.WriteString(_1)
		buffer.WriteString(_2)
		buffer.WriteString(_3)
		_ = buffer.String()
	}
}

func BenchmarkA4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var buffer bytes.Buffer
		buffer.Grow(len(_1) + len(_2) + len(_3))

		buffer.WriteString(_1)
		buffer.WriteString(_2)
		buffer.WriteString(_3)
		_ = buffer.String()
	}
}

func BenchmarkA5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buffer := acquireBytesBuffer()
		buffer.WriteString(_1)
		buffer.WriteString(_2)
		buffer.WriteString(_3)
		_ = buffer.String()
		poolBuffers.Put(buffer)
	}
}

func BenchmarkA6(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var buffer bytes.Buffer
		buffer.WriteString(_1)
		buffer.WriteString(_2)
		buffer.WriteString(_3)
		_ = buffer.String()
	}
}

func BenchmarkA7(b *testing.B) {
	for n := 0; n < b.N; n++ {
		builder := acquireStringsBuilder()
		builder.WriteString(_1)
		builder.WriteString(_2)
		builder.WriteString(_3)
		_ = builder.String()
		poolBuilders.Put(builder)
	}
}

func BenchmarkA8(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var builder strings.Builder
		builder.WriteString(_1)
		builder.WriteString(_2)
		builder.WriteString(_3)
		_ = builder.String()
	}
}
