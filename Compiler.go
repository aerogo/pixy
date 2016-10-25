package pixy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// ASTNode ...
type ASTNode struct {
	Line     string
	Children []*ASTNode
	Parent   *ASTNode
	Indent   int
}

// BuildAST returns a tree structure if you feed it with indentantion based source code.
func BuildAST(src string) *ASTNode {
	ast := new(ASTNode)
	ast.Indent = -1

	block := ast
	lastNode := ast

	lines := strings.Split(src, "\n")

	for _, line := range lines {
		// Ignore empty lines
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Indentation
		indent := 0
		for indent < len(line) {
			if line[indent] != '\t' {
				break
			}

			indent++
		}

		if indent != 0 {
			line = line[indent:]
		}

		node := new(ASTNode)
		node.Line = line
		node.Indent = indent

		if node.Indent == block.Indent+1 {
			// OK
		} else if node.Indent == block.Indent+2 {
			block = lastNode
		} else if node.Indent == block.Indent {
			block = block.Parent
		} else {
			panic("Invalid indentation")
		}

		node.Parent = block
		block.Children = append(block.Children, node)

		lastNode = node
	}

	return ast
}

// CompileFile compiles a Pixy template from fileIn and writes the
// resulting Go code to fileOut. It also returns the Go code as a string.
func CompileFile(fileIn string, fileOut string) string {
	srcBytes, _ := ioutil.ReadFile(fileIn)
	src := string(srcBytes)
	code := Compile(src)
	ioutil.WriteFile(fileOut, []byte(code), 0644)
	return code
}

// Compile compiles a Pixy template as a string and returns Go code.
func Compile(src string) string {
	ast := BuildAST(src)
	dir, _ := os.Getwd()
	packageName := filepath.Base(dir)
	code := compileChildren(ast)
	htmlPackageUsed := strings.Index(code, "html.EscapeString(") != -1
	packages := []string{"bytes"}

	if htmlPackageUsed {
		packages = append(packages, "html")
	}

	return "package " + packageName + "\n\nimport (\n\t\"" + strings.Join(packages, "\"\n\t\"") + "\"\n)\n\n" + code
}

// Compiles the children of a Pixy ASTNode.
func compileChildren(node *ASTNode) string {
	output := ""

	for _, child := range node.Children {
		code := compileNode(child)
		if len(code) > 0 {
			output += code + "\n"
		}
	}

	return strings.TrimSpace(output)
}

func write(s string) string {
	return "_b.WriteString(" + s + ")"
}

// Compiles a single Pixy ASTNode.
func compileNode(node *ASTNode) string {
	var keyword string

	for i, letter := range node.Line {
		// Function calls
		if i == 0 && unicode.IsLetter(letter) && unicode.IsUpper(letter) {
			return write(node.Line)
		}

		if len(keyword) == 0 && !unicode.IsLetter(letter) && !unicode.IsDigit(letter) {
			keyword = string([]rune(node.Line)[:i])
			fmt.Println("KEYWORD", keyword)
		}
	}

	if len(keyword) == 0 {
		keyword = node.Line
	}

	if keyword == "component" {
		functionBody := "var _b bytes.Buffer\n" + compileChildren(node) + "\nreturn _b.String()"
		lines := strings.Split(functionBody, "\n")

		if !strings.HasSuffix(node.Line, ")") {
			node.Line += "()"
		}

		comment := "// " + node.Line[len(keyword)+1:strings.Index(node.Line, "(")] + " component"
		return comment + "\nfunc " + node.Line[len("component "):] + " string {\n\t" + strings.Join(lines, "\n\t") + "\n}"
	}

	if keyword == "header" {
		return compileChildren(node)
	}

	if keyword == "h1" {
		var contents string

		// No contents?
		if node.Line == keyword {
			return write("\"<" + keyword + "></" + keyword + ">\"")
		}

		escapeInput := true
		equalIndex := len(keyword)

		if node.Line[equalIndex] == '!' {
			escapeInput = false
			equalIndex++
		}

		if node.Line[equalIndex] == '=' {
			contents = strings.TrimLeft(node.Line[equalIndex+1:], " ")

			code := write("\"<"+keyword+">\"") + "\n"

			if escapeInput {
				code += write("html.EscapeString("+contents+")") + "\n"
			} else {
				code += write(contents) + "\n"
			}

			code += write("\"</" + keyword + ">\"")
			return code
		}

		contents = strings.TrimLeft(node.Line[len(keyword):], " ")
		contents = strings.Replace(contents, "\"", "\\\"", -1)
		return write("\"<" + keyword + ">" + contents + "</" + keyword + ">\"")
		// return write("\"<h1>\" + html.EscapeString(\"" + contents + "\") + \"</h1>\"")
	}

	return "// Parse error: [" + node.Line + "]"
}
