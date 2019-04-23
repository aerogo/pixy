package pixy_test

import (
	"io/ioutil"
	"testing"

	"github.com/aerogo/pixy"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	src, _ := ioutil.ReadFile("examples/post-benchmark.pixy")
	code := string(src)

	components, err := pixy.Compile(code)
	assert.NoError(t, err)
	assert.NotNil(t, components)
	assert.Len(t, components, 1)
}

func TestCompileFile(t *testing.T) {
	components, err := pixy.CompileFile("examples/post-benchmark.pixy")
	assert.NoError(t, err)
	assert.NotNil(t, components)
	assert.Len(t, components, 1)
}

func BenchmarkCompile(b *testing.B) {
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