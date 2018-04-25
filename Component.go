package pixy

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
)

// Component represents a single, reusable template.
type Component struct {
	Name               string
	InterfaceCode      string
	ImplementationCode string
}

// Save writes the component to the given directory.
func (component *Component) Save(dirOut string) (interfaceFile string, implementationFile string) {
	// Write interface file
	interfaceFile = path.Join(dirOut, component.Name+".go")
	writeErr := ioutil.WriteFile(interfaceFile, []byte(component.InterfaceCode), 0644)

	if writeErr != nil {
		color.Red("Can't write to " + interfaceFile)
		color.Red(writeErr.Error())
	}

	// Write implementation file
	packageDirectory := path.Join(dirOut, "stream"+strings.ToLower(component.Name))
	os.MkdirAll(packageDirectory, 0777)

	implementationFile = path.Join(packageDirectory, component.Name+".go")
	writeErr = ioutil.WriteFile(implementationFile, []byte(component.ImplementationCode), 0644)

	if writeErr != nil {
		color.Red("Can't write to " + implementationFile)
		color.Red(writeErr.Error())
	}

	return interfaceFile, implementationFile
}
