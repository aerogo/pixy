package main

import (
	"io/ioutil"

	"github.com/aerogo/pixy"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.MemProfile).Stop()

	src, _ := ioutil.ReadFile("examples/post-benchmark.pixy")
	code := string(src)

	for i := 0; i < 100000000; i++ {
		_, err := pixy.Compile(code)

		if err != nil {
			panic(err)
		}
	}
}
