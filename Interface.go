package pixy

import "io"

// Compile compiles a Pixy template as a reader and returns a slice of components.
func Compile(reader io.Reader) ([]*Component, error) {
	return DefaultCompiler.Compile(reader)
}

// CompileString compiles a Pixy template as a byte slice and returns a slice of components.
func CompileString(src string) ([]*Component, error) {
	return DefaultCompiler.CompileString(src)
}

// CompileBytes compiles a Pixy template as a byte slice and returns a slice of components.
func CompileBytes(src []byte) ([]*Component, error) {
	return DefaultCompiler.CompileBytes(src)
}

// CompileFile compiles a Pixy template read from a file and returns a slice of components.
func CompileFile(fileIn string) ([]*Component, error) {
	return DefaultCompiler.CompileFile(fileIn)
}
