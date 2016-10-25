package pixy

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestCompiler(t *testing.T) {
	srcBytes, _ := ioutil.ReadFile("examples/hello.pixy")
	src := string(srcBytes)
	code := Compile(src)
	fmt.Println(code)
}

func BenchmarkCompiler(b *testing.B) {
	srcBytes, _ := ioutil.ReadFile("examples/hello.pixy")
	src := string(srcBytes)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Compile(src)
	}
}
