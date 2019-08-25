package pixy_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/aerogo/pixy"
	"github.com/akyoto/assert"
)

func TestCompile(t *testing.T) {
	file, err := os.Open("testdata/post-benchmark.pixy")
	assert.Nil(t, err)
	defer file.Close()

	components, err := pixy.Compile(file)
	assert.Nil(t, err)
	assert.NotNil(t, components)
	assert.Equal(t, len(components), 1)
}

func TestCompileBytes(t *testing.T) {
	code, _ := ioutil.ReadFile("testdata/post-benchmark.pixy")

	components, err := pixy.CompileBytes(code)
	assert.Nil(t, err)
	assert.NotNil(t, components)
	assert.Equal(t, len(components), 1)
}

func TestCompileString(t *testing.T) {
	src, _ := ioutil.ReadFile("testdata/post-benchmark.pixy")
	code := string(src)

	components, err := pixy.CompileString(code)
	assert.Nil(t, err)
	assert.NotNil(t, components)
	assert.Equal(t, len(components), 1)
}

func TestCompileFile(t *testing.T) {
	components, err := pixy.CompileFile("testdata/post-benchmark.pixy")
	assert.Nil(t, err)
	assert.NotNil(t, components)
	assert.Equal(t, len(components), 1)
}

func BenchmarkCompileString(b *testing.B) {
	src, _ := ioutil.ReadFile("testdata/post-benchmark.pixy")
	code := string(src)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := pixy.CompileString(code)

			if err != nil {
				b.Fail()
			}
		}
	})
}
