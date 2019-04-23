package pixy

import (
	"io/ioutil"
	"path"

	"github.com/akyoto/color"
)

// Component represents a single, reusable template.
type Component struct {
	Name string
	Code []byte
}

// Save writes the component to the given directory and returns the file path.
func (component *Component) Save(dirOut string) string {
	// Write interface file
	file := path.Join(dirOut, component.Name+".go")
	writeErr := ioutil.WriteFile(file, component.Code, 0644)

	if writeErr != nil {
		color.Red("Can't write to " + file)
		color.Red(writeErr.Error())
	}

	return file
}
