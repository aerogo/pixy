package main

import (
	"io/ioutil"

	"github.com/fatih/color"
)

// resulting Go code to fileOut. It also returns the Go code as a string.
func CompileFileAndSave(fileIn string, fileOut string) string {
	code := CompileFile(fileIn, true)
	writeErr := ioutil.WriteFile(fileOut, []byte(code), 0644)

	if writeErr != nil {
		color.Red("Can't write to " + fileOut)
	}

	return code
}

// CompileFile compiles a Pixy template from fileIn and returns the Go code as a string.
func CompileFile(fileIn string, includeHeader bool) string {
	srcBytes, readErr := ioutil.ReadFile(fileIn)

	if readErr != nil {
		color.Red("Can't read from " + fileIn)
		return ""
	}

	src := string(srcBytes)
	return Compile(src, includeHeader)
}

// Compile compiles a Pixy template as a string and returns Go code.
func Compile(src string, includeHeader bool) string {
	ast := BuildAST(src)
	code := compileChildren(ast)

	if includeHeader {
		return buildHeader(code) + code
	}

	return optimize(code)
}
