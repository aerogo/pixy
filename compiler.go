package main

import (
	"strings"
	"unicode"

	"github.com/fatih/color"
)

// PackageName contains the package name used in the generated .go files.
var PackageName = "components"

// Builds the file header.
func buildHeader(code string) string {
	return "package " + PackageName + "\n\n"
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

// Writes expression to the output.
func write(expression string) string {
	return "_b.WriteString(" + expression + ")\n"
}

// Writes s interpreted as a string (not an expression) to the output.
func writeString(s string) string {
	return write("\"" + s + "\"")
}

// Compiles a single ASTNode.
func compileNode(node *ASTNode) string {
	var keyword string

	for i, letter := range node.Line {
		// Function calls
		if i == 0 && unicode.IsLetter(letter) && unicode.IsUpper(letter) {
			if !strings.HasSuffix(node.Line, ")") {
				node.Line += "()"
			}

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

		if strings.HasSuffix(node.Line, "()") {
			color.Red(node.Line)
			color.Red("Components without parameters should not include parentheses in the definition.")
		}

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
