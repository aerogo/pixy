package pixy

// Compile compiles a Pixy template as a string and returns a slice of components.
func Compile(src string) ([]*Component, error) {
	return DefaultCompiler.Compile(src)
}

// CompileBytes compiles a Pixy template as a byte slice and returns a slice of components.
func CompileBytes(src []byte) ([]*Component, error) {
	return DefaultCompiler.CompileBytes(src)
}

// CompileFile compiles a Pixy template read from a file and returns a slice of components.
func CompileFile(fileIn string) ([]*Component, error) {
	return DefaultCompiler.CompileFile(fileIn)
}
