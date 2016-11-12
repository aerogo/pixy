package pixy

// PackageName contains the package name used in the generated .go files.
var PackageName = "components"

// Builds the file header.
func getFileHeader() string {
	return "package " + PackageName + "\n"
}

// Utility functions
func getUtilities() string {
	return getFileHeader() + `
import (
	"sync"
	"bytes"
)

var pool sync.Pool

func acquireBytesBuffer() *bytes.Buffer {
	var _b *bytes.Buffer
	obj := pool.Get()

	if obj == nil {
		return &bytes.Buffer{}
	}

	_b = obj.(*bytes.Buffer)
	_b.Reset()
	return _b
}
`
}
