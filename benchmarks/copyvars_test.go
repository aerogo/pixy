package pixy

import (
	"bytes"
	"sync"
	"testing"
)

var poolBuffers sync.Pool

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

func BenchmarkA1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_1 := "unknown string"
		_2 := "another string"
		_3 := "yet another string"

		_b := make([]byte, len(_1)+len(_2)+len(_3))
		_l := 0

		_l += copy(_b[_l:], _1)
		_l += copy(_b[_l:], _2)
		_l += copy(_b[_l:], _3)

		_ = string(_b)
	}
}

func BenchmarkA2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_s := [3]string{}
		_i := 0

		_s[_i] = "unknown string"
		_i++
		_s[_i] = "another string"
		_i++
		_s[_i] = "yet another string"
		_i++

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
		a := "unknown string"
		b := "another string"
		c := "yet another string"

		var buffer bytes.Buffer
		buffer.WriteString(a)
		buffer.WriteString(b)
		buffer.WriteString(c)
		_ = buffer.String()
	}
}

func BenchmarkA4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		a := "unknown string"
		b := "another string"
		c := "yet another string"

		var buffer bytes.Buffer
		buffer.Grow(len(a) + len(b) + len(c))

		buffer.WriteString(a)
		buffer.WriteString(b)
		buffer.WriteString(c)
		_ = buffer.String()
	}
}

func BenchmarkA5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		a := "unknown string"
		b := "another string"
		c := "yet another string"

		buffer := acquireBytesBuffer()
		buffer.WriteString(a)
		buffer.WriteString(b)
		buffer.WriteString(c)
		_ = buffer.String()
		poolBuffers.Put(buffer)
	}
}
