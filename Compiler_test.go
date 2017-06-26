package pixy

import (
	"io/ioutil"
	"testing"
)

func BenchmarkCompiler(b *testing.B) {
	src, _ := ioutil.ReadFile("examples/post-benchmark.pixy")
	code := string(src)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Compile(code)
		}
	})
}
