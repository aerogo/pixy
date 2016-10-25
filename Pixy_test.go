package pixy

import (
	"fmt"
	"io/ioutil"
	"testing"

	components "./examples"
)

func TestCompiler(t *testing.T) {
	code := CompileFile("examples/hello.pixy", "examples/hello.go")

	fmt.Println("--------------------------------------------------------------------")
	fmt.Println(code)
	fmt.Println("--------------------------------------------------------------------")
}

func TestExample(t *testing.T) {
	html := components.Hello("Pixy")
	ioutil.WriteFile("examples/hello.html", []byte(html), 0644)

	fmt.Println(html)
	fmt.Println("--------------------------------------------------------------------")
}

func BenchmarkCompiler(b *testing.B) {
	srcBytes, _ := ioutil.ReadFile("examples/hello.pixy")
	src := string(srcBytes)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Compile(src)
	}
}
