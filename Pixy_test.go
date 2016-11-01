package main

import "testing"

func TestCompiler(t *testing.T) {
	main()
}

func BenchmarkCompiler(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		main()
	}
}
