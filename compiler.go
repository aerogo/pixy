package main

import (
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

// PackageName contains the package name used in the generated .go files.
var PackageName = "components"

// ASTNode ...
type ASTNode struct {
	Line     string
	Children []*ASTNode
	Parent   *ASTNode
	Indent   int
}

var compactCode *regexp.Regexp

// init
func init() {
	compactCode = regexp.MustCompile("\\n{2,}")
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

// CompileFileAndSave compiles a Pixy template from fileIn and writes the
// resulting Go code to fileOut. It also returns the Go code as a string.
func CompileFileAndSave(fileIn string, fileOut string) string {
	code := CompileFile(fileIn, true)
	writeErr := ioutil.WriteFile(fileOut, []byte(code), 0644)

	if writeErr != nil {
		color.Red("Can't write to " + fileOut)
	}

	return code
}

// CompileFile compiles a Pixy template from fileIn and returns the Go code as a string.
func CompileFile(fileIn string, includeHeader bool) string {
	srcBytes, readErr := ioutil.ReadFile(fileIn)

	if readErr != nil {
		color.Red("Can't read from " + fileIn)
		return ""
	}

	src := string(srcBytes)
	return Compile(src, includeHeader)
}

// Compile compiles a Pixy template as a string and returns Go code.
func Compile(src string, includeHeader bool) string {
	ast := BuildAST(src)
	code := compileChildren(ast)

	if includeHeader {
		return buildHeader(code) + code
	}

	return optimize(code)
}

// buildHeader ...
func buildHeader(code string) string {
	return "package " + PackageName + "\n\n"
}

const writeStringCall = "_b.WriteString(\""

// optimize combines multiple WriteString calls to one.
func optimize(code string) string {
	lines := strings.Split(code, "\n")
	lastString := ""

	for index, line := range lines {
		pos := strings.Index(line, writeStringCall)

		if pos != -1 {
			lastString += line[pos+len(writeStringCall) : len(line)-2]
			lines[index] = ""
			continue
		}

		if len(lastString) > 0 {
			lines[index] = "\t" + writeStringCall + lastString + "\")\n" + line
			lastString = ""
		}
	}

	return compactCode.ReplaceAllString(strings.Join(lines, "\n"), "\n")
}

// Compiles the children of a Pixy ASTNode.
func compileChildren(node *ASTNode) string {
	output := ""

	for _, child := range node.Children {
		code := strings.TrimSpace(compileNode(child))
		if len(code) > 0 {
			output += code + "\n"
		}
	}

	return output
}

func write(s string) string {
	return "_b.WriteString(" + s + ")\n"
}

func writeString(s string) string {
	return write("\"" + s + "\"")
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
		}
	}

	if len(keyword) == 0 {
		keyword = node.Line
	}

	if keyword == "component" {
		functionBody := "var _b bytes.Buffer\n" + compileChildren(node) + "return _b.String()"
		lines := strings.Split(functionBody, "\n")

		if !strings.HasSuffix(node.Line, ")") {
			node.Line += "()"
		}

		comment := "// " + node.Line[len(keyword)+1:strings.Index(node.Line, "(")] + " component"
		return comment + "\nfunc " + node.Line[len("component "):] + " string {\n\t" + strings.Join(lines, "\n\t") + "\n}"
	}

	var contents string
	attributes := make(map[string]string)

	tag := func() string {
		if len(attributes) == 0 {
			return writeString("<" + keyword + ">")
		}

		code := writeString("<" + keyword + " ")
		for key, value := range attributes {
			code += writeString(key + "=")
			code += write(value)
		}
		code += writeString(">")
		return code
	}

	// No contents?
	if node.Line == keyword {
		code := ""

		if keyword == "html" {
			code += writeString("<!DOCTYPE html>")
		}

		code += tag()
		code += compileChildren(node)
		code += writeString("</" + keyword + ">")
		return code
	}

	escapeInput := true
	cursor := len(keyword)

	if node.Line[cursor] == '#' {
		cursor++
		start := cursor
		analyze := node.Line[cursor:]
		for index, letter := range analyze {
			if !unicode.IsLetter(letter) && !unicode.IsDigit(letter) && letter != '-' {
				cursor += index
				id := node.Line[start:cursor]
				attributes["id"] = "\"" + id + "\""
				break
			}
		}
	}

	if node.Line[cursor] == '!' {
		escapeInput = false
		cursor++
	}

	if node.Line[cursor] == '=' {
		contents = strings.TrimLeft(node.Line[cursor+1:], " ")

		code := tag()

		if escapeInput {
			code += write("html.EscapeString(" + contents + ")")
		} else {
			code += write(contents)
		}

		code += compileChildren(node)
		code += writeString("</" + keyword + ">")
		return code
	}

	contents = strings.TrimLeft(node.Line[len(keyword):], " ")
	contents = strings.Replace(contents, "\"", "\\\"", -1)
	code := tag()
	code += writeString(contents)
	code += compileChildren(node)
	code += writeString("</" + keyword + ">")
	return code
	// return write("\"<h1>\" + html.EscapeString(\"" + contents + "\") + \"</h1>\"")

	// return "// Parse error: [" + node.Line + "]"
}
