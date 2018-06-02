package pixy

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const testString = "xx"

var testBytes = []byte(testString)

// BenchmarkJoin ...
func BenchmarkJoin(b *testing.B) {
	c := make([]string, b.N)

	for n := 0; n < b.N; n++ {
		c[n] = testString
	}

	str := strings.Join(c, "")

	if len(str) != b.N*len(testString) {
		fmt.Println("Error BenchmarkJoin")
	}
}

// BenchmarkWrite ...
func BenchmarkWrite(b *testing.B) {
	var buffer bytes.Buffer

	for n := 0; n < b.N; n++ {
		buffer.Write(testBytes)
	}

	str := buffer.String()

	if len(str) != b.N*len(testString) {
		fmt.Println("Error BenchmarkWrite")
	}
}

// BenchmarkWriteString ...
func BenchmarkWriteString(b *testing.B) {
	var buffer bytes.Buffer

	for n := 0; n < b.N; n++ {
		buffer.WriteString(testString)
	}

	str := buffer.String()

	if len(str) != b.N*len(testString) {
		fmt.Println("Error BenchmarkWriteString")
	}
}

// BenchmarkWriteStringKL ...
func BenchmarkWriteStringKL(b *testing.B) {
	var buffer bytes.Buffer
	buffer.Grow(b.N)

	for n := 0; n < b.N; n++ {
		buffer.WriteString(testString)
	}

	str := buffer.String()

	if len(str) != b.N*len(testString) {
		fmt.Println("Error BenchmarkWriteStringKL")
	}
}

// BenchmarkCopyUL ...
func BenchmarkCopyUL(b *testing.B) {
	s := make([]string, b.N)
	_l := 0

	for n := 0; n < b.N; n++ {
		s[n] = testString
		_l += len(testString)
	}

	_b := make([]byte, _l)
	_c := 0

	for i := 0; i < b.N; i++ {
		_c += copy(_b[_c:], s[i])
	}

	str := string(_b)

	if len(str) != b.N*len(testString) {
		fmt.Println("Error BenchmarkCopyUL")
	}
}

// BenchmarkCopyStringKL ...
func BenchmarkCopyStringKL(b *testing.B) {
	_l := len(testString) * b.N
	_b := make([]byte, _l)
	_c := 0

	for n := 0; n < b.N; n++ {
		_c += copy(_b[_c:], testString)
	}

	str := string(_b)

	if len(str) != b.N*len(testString) {
		fmt.Println("Error BenchmarkCopyStringKL")
	}
}

// BenchmarkCopyBytesKL ...
func BenchmarkCopyBytesKL(b *testing.B) {
	_l := len(testString) * b.N
	_b := make([]byte, _l)
	_c := 0

	for n := 0; n < b.N; n++ {
		_c += copy(_b[_c:], testBytes)
	}

	str := string(_b)

	if len(str) != b.N*len(testString) {
		fmt.Println("Error BenchmarkCopyBytesKL")
	}
}

// BenchmarkMix ...
func BenchmarkMix(b *testing.B) {
	var buffer bytes.Buffer
	iterations := b.N

	if iterations == 1 {
		iterations = 2
	}

	half := iterations / 2
	_l := len(testString) * half
	_b := make([]byte, _l)
	_c := 0

	for n := 0; n < half; n++ {
		buffer.WriteString(testString)
	}

	for n := 0; n < half; n++ {
		_c += copy(_b[_c:], testBytes)
	}

	buffer.Write(_b)

	str := buffer.String()

	if len(str) != iterations*len(testString) {
		fmt.Println("Error BenchmarkMix")
	}
}
