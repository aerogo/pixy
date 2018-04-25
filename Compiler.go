package pixy

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/aerogo/codetree"
	"github.com/fatih/color"
)

// DefaultCompiler is the default compiler used by the interface.
var DefaultCompiler = NewCompiler("components")

// Compiler is a Pixy file compiler.
type Compiler struct {
	// PackageName contains the package name used in the generated .go files.
	PackageName string
}

// NewCompiler constructs a new Pixy compiler.
func NewCompiler(packageName string) *Compiler {
	return &Compiler{
		PackageName: packageName,
	}
}

// Compile compiles a Pixy template as a string and returns a slice of components.
func (compiler *Compiler) Compile(src string) ([]*Component, error) {
	tree, err := codetree.New(src)
	defer tree.Close()

	if err != nil {
		return nil, err
	}

	components := []*Component{}

	for _, node := range tree.Children {
		// Disallow tags on the top level
		if !strings.HasPrefix(node.Line, "component ") && !strings.HasPrefix(node.Line, "//") {
			color.Yellow(node.Line)
			color.Red("Only 'component' definitions are allowed on the top level.")
			continue
		}

		// Signature contains the signature of the component without the preceding keyword.
		signature := node.Line[len("component "):]

		// Any signature that ends with empty parentheses should be rewritten to not include them.
		if strings.HasSuffix(signature, "()") {
			color.Yellow(signature)
			color.Red("Components without definition should not include parentheses in the definition.")
		}

		// Add parentheses to empty parameter lists
		if !strings.HasSuffix(signature, ")") {
			signature += "()"
		}

		// Get the necessary info from the component signature
		componentName := signature[:strings.Index(signature, "(")]
		componentParameters := signature[len(componentName)+1 : len(signature)-1]
		parameterNames := extractParameterNames(componentParameters)

		// streamFunctionCall contains the function call for the streaming version.
		streamFunctionCall := "stream" + strings.ToLower(componentName) + ".Stream" + componentName + "(_b"

		if len(parameterNames) > 0 {
			streamFunctionCall += ", " + strings.Join(parameterNames, ", ")
		}

		streamFunctionCall += ")"

		// Generate a comment line so that the linter won't complain
		comment := "// " + componentName + " component"

		// Stream function body
		streamFunctionBody := compileChildren(node)
		streamFunctionBody = strings.Replace(streamFunctionBody, "\n", "\n\t", -1)
		optimizedStreamFunctionBody, inlined := optimize(streamFunctionBody)

		// Normal function body
		functionBody := ""

		if inlined != "" {
			functionBody = strings.TrimSpace(inlined)
		} else {
			functionBody = "_b := acquireBytesBuffer()\n" + streamFunctionCall + "\npool.Put(_b)\nreturn _b.String()"
			functionBody = strings.Replace(functionBody, "\n", "\n\t", -1)
		}

		// Build the component code
		interfaceCode := acquireBytesBuffer()
		implementationCode := acquireBytesBuffer()

		// Normal function
		interfaceCode.WriteString(compiler.GetFileHeader())
		interfaceCode.WriteString(comment)
		interfaceCode.WriteString("\nfunc ")
		interfaceCode.WriteString(signature)
		interfaceCode.WriteString(" string {\n\t")
		interfaceCode.WriteString(functionBody)
		interfaceCode.WriteString("\n}")

		// Stream function
		implementationCode.WriteString("package stream")
		implementationCode.WriteString(strings.ToLower(componentName))
		implementationCode.WriteByte('\n')
		implementationCode.WriteByte('\n')
		implementationCode.WriteString("\nfunc Stream")
		implementationCode.WriteString(strings.Replace(signature, "(", "(_b *bytes.Buffer, ", 1))
		implementationCode.WriteString(" {")
		implementationCode.WriteString(optimizedStreamFunctionBody)
		implementationCode.WriteString("}")

		// Add the compiled component to the return values
		components = append(components, &Component{
			Name:               componentName,
			InterfaceCode:      interfaceCode.String(),
			ImplementationCode: implementationCode.String(),
		})

		// Allow the byte buffer to be re-used
		pool.Put(interfaceCode)
	}

	return components, nil
}

// CompileBytes compiles a Pixy template as a byte slice and returns a slice of components.
func (compiler *Compiler) CompileBytes(src []byte) ([]*Component, error) {
	return compiler.Compile(string(src))
}

// CompileFile compiles a Pixy template read from a file and returns a slice of components.
func (compiler *Compiler) CompileFile(fileIn string) ([]*Component, error) {
	src, err := ioutil.ReadFile(fileIn)

	if err != nil {
		return nil, errors.New("Can't read from " + fileIn + "\n" + err.Error())
	}

	return compiler.CompileBytes(src)
}

// CompileFileAndSaveIn compiles a Pixy template from fileIn
// and writes the resulting components to dirOut.
func (compiler *Compiler) CompileFileAndSaveIn(fileIn string, dirOut string) ([]*Component, []string, error) {
	components, err := compiler.CompileFile(fileIn)
	files := make([]string, len(components)*2, len(components)*2)

	for index, component := range components {
		files[index*2], files[index*2+1] = component.Save(dirOut)
	}

	return components, files, err
}

// GetFileHeader returns the file header.
func (compiler *Compiler) GetFileHeader() string {
	return "package " + compiler.PackageName + "\n"
}

// GetUtilities returns the file header and utility functions
// that are available for components.
func (compiler *Compiler) GetUtilities() string {
	return compiler.GetFileHeader() + `
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

// SaveUtilities adds the file with required function definitions to the directory.
func (compiler *Compiler) SaveUtilities(filePath string) {
	ioutil.WriteFile(filePath, []byte(compiler.GetUtilities()), 0644)
}
