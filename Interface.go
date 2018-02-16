package pixy

import (
	"io/ioutil"
	"strings"

	"github.com/aerogo/codetree"
	"github.com/fatih/color"
)

// CompileFileAndSaveIn compiles a Pixy template from fileIn
// and writes the resulting components to dirOut.
func CompileFileAndSaveIn(fileIn string, dirOut string) []*Component {
	srcBytes, readErr := ioutil.ReadFile(fileIn)

	if readErr != nil {
		color.Red("Can't read from " + fileIn)
		return nil
	}

	src := string(srcBytes)
	components := Compile(src)

	for _, component := range components {
		component.Save(dirOut)
	}

	return components
}

// SaveUtilities adds the file $.go with required function definitions to the directory.
func SaveUtilities(filePath string) {
	ioutil.WriteFile(filePath, []byte(getUtilities()), 0644)
}

// Compile compiles a Pixy template as a string and returns a slice of components.
func Compile(src string) []*Component {
	tree, err := codetree.New(src)
	defer tree.Close()

	if err != nil {
		panic(err)
	}

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

		componentName := definition[:strings.Index(definition, "(")]
		streamFunctionBody := compileChildren(node)
		functionBody := "_b := acquireBytesBuffer()\n" + streamFunctionBody + "pool.Put(_b)\nreturn _b.String()"
		functionBody = strings.Replace(functionBody, "\n", "\n\t", -1)
		streamFunctionBody = strings.Replace(streamFunctionBody, "\n", "\n\t", -1)
		comment := "// " + componentName + " component"

		componentCode := acquireBytesBuffer()
		componentCode.WriteString(getFileHeader())

		// Normal function
		componentCode.WriteString(comment)
		componentCode.WriteString("\nfunc ")
		componentCode.WriteString(definition)
		componentCode.WriteString(" string {\n\t")
		componentCode.WriteString(optimize(functionBody))
		componentCode.WriteString("\n}")

		// Stream function
		componentCode.WriteByte('\n')
		componentCode.WriteString("\nfunc stream")
		componentCode.WriteString(strings.Replace(definition, "(", "(_b *bytes.Buffer, ", 1))
		componentCode.WriteString(" {")
		componentCode.WriteString(optimize(streamFunctionBody))
		componentCode.WriteString("}")

		components = append(components, &Component{
			Name: componentName,
			Code: componentCode.String(),
		})

		pool.Put(componentCode)
	}

	return components
}
