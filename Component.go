package pixy

import (
	"io/ioutil"
	"path"

	"github.com/fatih/color"
)

// Component represents a single, reusable template.
type Component struct {
	Name string
	Code string
}

// Save writes the component to the given directory.
func (component *Component) Save(dirOut string) {
	fileOut := path.Join(dirOut, component.Name+".go")
	writeErr := ioutil.WriteFile(fileOut, []byte(component.Code), 0644)

	if writeErr != nil {
		color.Red("Can't write to " + fileOut)
		color.Red(writeErr.Error())
	}

	// Run goimports
	goimports(fileOut)
}
