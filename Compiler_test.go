package pixy_test

import (
	"io/ioutil"
	"testing"

	"github.com/aerogo/pixy"
)

func BenchmarkCompiler(b *testing.B) {
	src, _ := ioutil.ReadFile("examples/post-benchmark.pixy")
	code := string(src)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pixy.Compile(code)
		}
	})
}
