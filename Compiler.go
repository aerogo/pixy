package pixy

import (
	"strings"
	"unicode"

	"github.com/fatih/color"
)

// Compiles the children of a Pixy CodeTree.
func compileChildren(node *CodeTree) string {
	output := ""

	for _, child := range node.Children {
		code := strings.TrimSpace(compileNode(child))
		if len(code) > 0 {
			if strings.HasPrefix(code, "else {") {
				output = strings.TrimRight(output, "\n") + code + "\n"
			} else {
				output += code + "\n"
			}
		}
	}

	return output
}

// Writes expression to the output.
func write(expression string) string {
	if strings.HasPrefix(expression, "'") {
		color.Red("Strings must use \" instead of '")
		return ""
	}

	return "_b.WriteString(" + expression + ")\n"
}

// Writes s interpreted as a string (not an expression) to the output.
func writeString(s string) string {
	return write("\"" + s + "\"")
}

// isString
func isString(code string) bool {
	// TODO: Fix this
	return strings.HasPrefix(code, "\"") && strings.HasSuffix(code, "\"")
}

// Compiles a single CodeTree.
func compileNode(node *CodeTree) string {
	var keyword string

	if node.Line[0] == '#' || node.Line[0] == '.' {
		node.Line = "div" + node.Line
	}

	for i, letter := range node.Line {
		// Function calls
		if i == 0 && unicode.IsLetter(letter) && unicode.IsUpper(letter) {
			if !strings.HasSuffix(node.Line, ")") {
				node.Line += "()"
			}

			return write(node.Line)
		}

		// Comments
		if i == 1 && node.Line[0] == '/' && node.Line[1] == '/' {
			return ""
		}

		// Find keyword
		if len(keyword) == 0 && !unicode.IsLetter(letter) && !unicode.IsDigit(letter) && letter != '-' {
			keyword = string([]rune(node.Line)[:i])
		}
	}

	// Keyword takes full line
	if len(keyword) == 0 {
		keyword = node.Line
	}

	// Flow control
	if keyword == "if" || keyword == "else" || keyword == "for" {
		return node.Line + " {\n" + compileChildren(node) + "}"
	}

	// Each is just syntax sugar
	if keyword == "each" {
		// TODO: This is a just quick prototype implementation and not correct at all
		return strings.Replace(strings.Replace(node.Line, "each", "for _, ", 1), " in ", " := range ", 1) + " {\n" + compileChildren(node) + "}"
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
				code += write("html.EscapeString(fmt.Sprint(" + value + "))")
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

	endTag := func() string {
		if selfClosingTags[keyword] != true {
			return writeString("</" + keyword + ">")
		}

		return ""
	}

	// No contents?
	if node.Line == keyword {
		code := ""

		if keyword == "html" {
			code = writeString("<!DOCTYPE html>")
		}

		code += tag()
		code += compileChildren(node)
		code += endTag()
		return code
	}

	escapeInput := true
	cursor := len(keyword)

	expect := func(expected byte, callback func(int, string)) bool {
		if cursor >= len(node.Line) {
			return false
		}

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
		endFound := false

		for index, letter := range remaining {
			if !unicode.IsLetter(letter) && !unicode.IsDigit(letter) && letter != '-' {
				cursor += index
				id := node.Line[start:cursor]
				attributes["id"] = "\"" + id + "\""
				endFound = true
				break
			}
		}

		if !endFound {
			cursor = len(node.Line)
			id := node.Line[start:cursor]
			attributes["id"] = "\"" + id + "\""
		}
	})

	// Classes
	var classes []string
	for expect('.', func(start int, remaining string) {
		endFound := false

		for index, letter := range remaining {
			if !unicode.IsLetter(letter) && !unicode.IsDigit(letter) && letter != '-' {
				cursor += index
				name := node.Line[start:cursor]
				classes = append(classes, name)
				endFound = true
				break
			}
		}

		if !endFound {
			cursor = len(node.Line)
			name := node.Line[start:cursor]
			classes = append(classes, name)
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
				// fmt.Println("NAME", attributeName)
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

					if strings.HasPrefix(attributeValue, "'") {
						color.Yellow(attributeValue)
						color.Red("Strings must use \" instead of '")
						attributeValue = ""
					}

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

	if cursor < len(node.Line) {
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
				code += write("html.EscapeString(fmt.Sprint(" + contents + "))")
			} else {
				code += write(contents)
			}

			code += compileChildren(node)
			code += endTag()
			return code
		}

		contents = node.Line[cursor+1:]
		contents = strings.Replace(contents, "\"", "\\\"", -1)
	} else {
		contents = ""
	}

	code := tag()
	code += writeString(contents)
	code += compileChildren(node)
	code += endTag()
	return code
}
