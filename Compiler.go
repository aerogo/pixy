package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

// PackageName contains the package name used in the generated .go files.
var PackageName = "components"

// Builds the file header.
func getHeader() string {
	header := "package " + PackageName + "\n"
	header += `
type renderer struct{}

// Render methods allow you to render your components
var Render renderer

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
	return header
}

// Compiles the children of a Pixy CodeTree.
func compileChildren(node *CodeTree) string {
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

func isString(code string) bool {
	// TODO: Fix this
	return strings.HasPrefix(code, "\"") && strings.HasSuffix(code, "\"")
}

type ignoreReader struct {
	inString          bool
	inCharacterString bool
	inParentheses     int
	escape            bool
}

func (r *ignoreReader) canIgnore(letter rune) bool {
	if letter == '\\' && !r.escape {
		r.escape = true
		return true
	}

	defer func() {
		r.escape = false
	}()

	if letter == '"' && !r.escape {
		r.inString = !r.inString
		return true
	}

	if r.inString {
		return true
	}

	if letter == '\'' && !r.escape {
		r.inCharacterString = !r.inCharacterString
		return true
	}

	if r.inCharacterString {
		return true
	}

	if letter == '(' || letter == '[' || letter == '{' {
		r.inParentheses++
		return true
	}

	if letter == ')' || letter == ']' || letter == '}' {
		r.inParentheses--

		if r.inParentheses == 0 {
			return true
		}
	}

	if r.inParentheses > 0 {
		return true
	}

	return false
}

// Compiles a single CodeTree.
func compileNode(node *CodeTree) string {
	var keyword string

	for i, letter := range node.Line {
		// Function calls
		if i == 0 && unicode.IsLetter(letter) && unicode.IsUpper(letter) {
			if !strings.HasSuffix(node.Line, ")") {
				node.Line += "()"
			}

			return write("Render." + node.Line)
		}

		// Comments
		if i == 1 && node.Line[0] == '/' && node.Line[1] == '/' {
			return ""
		}

		// Find keyword
		if len(keyword) == 0 && !unicode.IsLetter(letter) && !unicode.IsDigit(letter) {
			keyword = string([]rune(node.Line)[:i])
		}
	}

	if len(keyword) == 0 {
		keyword = node.Line
	}

	if keyword == "component" {
		functionBody := "_b := acquireBytesBuffer()\n" + compileChildren(node) + "pool.Put(_b)\nreturn _b.String()"
		lines := strings.Split(functionBody, "\n")

		if strings.HasSuffix(node.Line, "()") {
			color.Yellow(node.Line)
			color.Red("Components without parameters should not include parentheses in the definition.")
		}

		if !strings.HasSuffix(node.Line, ")") {
			node.Line += "()"
		}

		comment := "// " + node.Line[len(keyword)+1:strings.Index(node.Line, "(")] + " component"
		return comment + "\nfunc (r *renderer) " + node.Line[len("component "):] + " string {\n\t" + strings.Join(lines, "\n\t") + "\n}"
	}

	// Disallow tags on the top level
	if node.Indent == 0 {
		color.Yellow(node.Line)
		color.Red("Only 'component' definitions are allowed on the top level.")
		return ""
	}

	var contents string
	attributes := make(map[string]string)

	tag := func() string {
		numAttributes := len(attributes)

		if numAttributes == 0 {
			return writeString("<" + keyword + ">")
		}

		code := writeString("<" + keyword + " ")
		count := 1

		for key, value := range attributes {
			code += writeString(key + "='")

			if isString(value) {
				code += write(strings.Replace(value, "'", "\\\\'", -1))
			} else {
				code += write("html.EscapeString(" + value + ")")
			}

			if count == numAttributes {
				code += writeString("'")
			} else {
				code += writeString("' ")
			}

			count++
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

	expect := func(expected byte, callback func(int, string)) bool {
		char := node.Line[cursor]

		if char == expected {
			cursor++

			if callback != nil {
				start := cursor
				remaining := node.Line[cursor:]
				callback(start, remaining)
			}

			return true
		}

		return false
	}

	// ID
	expect('#', func(start int, remaining string) {
		for index, letter := range remaining {
			if !unicode.IsLetter(letter) && !unicode.IsDigit(letter) && letter != '-' {
				cursor += index
				id := node.Line[start:cursor]
				attributes["id"] = "\"" + id + "\""
				break
			}
		}
	})

	// Classes
	var classes []string
	for expect('.', func(start int, remaining string) {
		for index, letter := range remaining {
			if !unicode.IsLetter(letter) && !unicode.IsDigit(letter) && letter != '-' {
				cursor += index
				name := node.Line[start:cursor]
				classes = append(classes, name)
				break
			}
		}
	}) {
		// Empty loop
	}

	readOneAttribute := func(start int, remaining string) bool {
		for node.Line[cursor] == ' ' {
			cursor++
		}

		remaining = node.Line[cursor:]
		start = cursor

		var attributeName string

		for index, letter := range remaining {
			if !unicode.IsLetter(letter) && letter != '-' {
				cursor += index
				attributeName = node.Line[start:cursor]
				fmt.Println("NAME", attributeName)
				break
			}
		}

		char := node.Line[cursor]

		if char == '=' {
			cursor++
			start = cursor
			remaining = node.Line[cursor:]

			var ignore ignoreReader
			for index, letter := range remaining {
				if ignore.canIgnore(letter) {
					continue
				}

				if letter == ',' || letter == ')' {
					cursor += index
					attributeValue := node.Line[start:cursor]
					fmt.Println("VALUE", attributeValue)
					attributes[attributeName] = attributeValue
					cursor++

					if letter == ',' {
						return true
					}

					return false
				}
			}
		}

		return false
	}

	// Attributes
	expect('(', func(start int, remaining string) {
		for readOneAttribute(start, remaining) != false {
			start = cursor
			remaining = node.Line[cursor:]
		}
	})

	if len(classes) > 0 {
		attributes["class"] = "\"" + strings.Join(classes, " ") + "\""
	}

	// Bypass HTML escaping
	if node.Line[cursor] == '!' {
		escapeInput = false
		cursor++
	}

	// Expressions
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

	contents = strings.TrimLeft(node.Line[cursor:], " ")
	contents = strings.Replace(contents, "\"", "\\\"", -1)
	code := tag()
	code += writeString(contents)
	code += compileChildren(node)
	code += writeString("</" + keyword + ">")
	return code
}
