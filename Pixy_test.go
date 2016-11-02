package main

import (
	"os"
	"testing"
)

func TestCompiler(t *testing.T) {
	os.Remove("❖.go")

	main()
}

func BenchmarkCompiler(b *testing.B) {
	os.Remove("❖.go")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		main()
	}
}
