package pixy

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aerogo/codetree"
	"github.com/akyoto/color"
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
func (compiler *Compiler) Compile(reader io.Reader) ([]*Component, error) {
	tree, err := codetree.New(reader)

	if err != nil {
		return nil, err
	}

	defer tree.Close()
	components := []*Component{}

	for _, node := range tree.Children {
		// Ignore comments
		if strings.HasPrefix(node.Line, "//") {
			continue
		}

		// Disallow tags on the top level
		if !strings.HasPrefix(node.Line, "component ") {
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
		streamFunctionCall := "stream" + componentName + "(_b"

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
			functionBody = "_b := acquireStringsBuilder()\n" + streamFunctionCall + "\n_pool.Put(_b)\nreturn _b.String()"
			functionBody = strings.Replace(functionBody, "\n", "\n\t", -1)
		}

		// Build the component code
		code := acquireStringsBuilder()

		// Normal function
		code.WriteString(compiler.GetFileHeader())
		code.WriteString(comment)
		code.WriteString("\nfunc ")
		code.WriteString(signature)
		code.WriteString(" string {\n\t")
		code.WriteString(functionBody)
		code.WriteString("\n}")

		// Stream function
		code.WriteByte('\n')
		code.WriteByte('\n')
		code.WriteString("func stream")
		code.WriteString(strings.Replace(signature, "(", "(_b *strings.Builder, ", 1))
		code.WriteString(" {")
		code.WriteString(optimizedStreamFunctionBody)
		code.WriteString("}")

		// Add the compiled component to the return values
		components = append(components, &Component{
			Name: componentName,
			Code: code.String(),
		})

		// Allow the byte buffer to be re-used
		pool.Put(code)
	}

	return components, nil
}

// CompileBytes compiles a Pixy template as a byte slice and returns a slice of components.
func (compiler *Compiler) CompileBytes(src []byte) ([]*Component, error) {
	return compiler.Compile(bytes.NewReader(src))
}

// CompileString compiles a Pixy template as a string and returns a slice of components.
func (compiler *Compiler) CompileString(src string) ([]*Component, error) {
	return compiler.Compile(strings.NewReader(src))
}

// CompileFile compiles a Pixy template read from a file and returns a slice of components.
func (compiler *Compiler) CompileFile(fileIn string) ([]*Component, error) {
	reader, err := os.Open(fileIn)

	if err != nil {
		return nil, errors.New("Can't read from " + fileIn + "\n" + err.Error())
	}

	return compiler.Compile(reader)
}

// GetFileHeader returns the file header.
func (compiler *Compiler) GetFileHeader() string {
	return "package " + compiler.PackageName + "\n\n"
}

// GetUtilities returns the file header and utility functions
// that are available for components.
func (compiler *Compiler) GetUtilities() string {
	return compiler.GetFileHeader() + `
import (
	"strings"
	"sync"
)

var _pool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

func acquireStringsBuilder() *strings.Builder {
	builder := _pool.Get().(*strings.Builder)
	builder.Reset()
	return builder
}
`
}

// SaveUtilities adds the file with required function definitions to the directory.
func (compiler *Compiler) SaveUtilities(filePath string) error {
	return ioutil.WriteFile(filePath, []byte(compiler.GetUtilities()), 0644)
}
