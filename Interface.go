package pixy

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/fatih/color"
)

// Component represents a single, reusable template.
type Component struct {
	Name string
	Code string
}

// CompileFileAndSaveIn compiles a Pixy template from fileIn
// and writes the resulting components to dirOut.
func CompileFileAndSaveIn(fileIn string, dirOut string) {
	srcBytes, readErr := ioutil.ReadFile(fileIn)

	if readErr != nil {
		color.Red("Can't read from " + fileIn)
		return
	}

	src := string(srcBytes)
	components := Compile(src)

	for _, component := range components {
		fileOut := path.Join(dirOut, component.Name, ".go")
		writeErr := ioutil.WriteFile(fileOut, []byte(src), 0644)

		if writeErr != nil {
			color.Red("Can't write to " + fileOut)
		}

		// Run goimports
		goimports(fileOut)
	}
}

// Compile compiles a Pixy template as a string and returns a slice of components.
func Compile(src string) []*Component {
	tree := NewCodeTree(src)
	components := []*Component{}

	for _, node := range tree.Children {
		// Disallow tags on the top level
		if !strings.HasPrefix(node.Line, "component ") {
			color.Yellow(node.Line)
			color.Red("Only 'component' definitions are allowed on the top level.")
			continue
		}

		definition := node.Line[len("component "):]

		if strings.HasSuffix(definition, "()") {
			color.Yellow(definition)
			color.Red("Components without parameters should not include parentheses in the definition.")
		}

		if !strings.HasSuffix(definition, ")") {
			definition += "()"
		}

		componentName := definition[:strings.Index(node.Line, "(")]
		functionBody := "_b := acquireBytesBuffer()\n" + compileChildren(node) + "pool.Put(_b)\nreturn _b.String()"
		lines := strings.Split(functionBody, "\n")
		comment := "// " + componentName + " component"
		componentCode := getFileHeader()
		componentCode += comment + "\nfunc (r *renderer) " + definition + " string {\n\t" + strings.Join(lines, "\n\t") + "\n}"
		componentCode = optimize(componentCode)

		components = append(components, &Component{
			Name: componentName,
			Code: componentCode,
		})
	}

	return components
}
