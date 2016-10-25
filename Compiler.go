package pixy

import (
	"io/ioutil"
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

// CompileFile compiles a Pixy template to Go code.
func CompileFile(fileIn string, fileOut string) {
	srcBytes, _ := ioutil.ReadFile(fileIn)
	src := string(srcBytes)
	code := Compile(src)
	ioutil.WriteFile(fileOut, []byte(code), 0644)
}

// Compile compiles a Pixy template as a string and returns Go code.
func Compile(src string) string {
	ast := BuildAST(src)
	return "package main\n\n" + compileChildren(ast)
}

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

// Compiles a single Pixy ASTNode.
func compileNode(node *ASTNode) string {
	for _, firstLetter := range node.Line {
		if unicode.IsLetter(firstLetter) && unicode.IsUpper(firstLetter) {
			return node.Line
		}

		break
	}

	if strings.HasPrefix(node.Line, "component ") {
		functionBody := compileChildren(node)
		lines := strings.Split(functionBody, "\n")
		return "func " + node.Line[len("component "):] + " {\n\t" + strings.Join(lines, "\n\t") + "\n}"
	}

	if node.Line == "img" {
		return "<img src=''>"
	}

	return ""
}
