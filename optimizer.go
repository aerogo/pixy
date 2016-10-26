package main

import (
	"regexp"
	"strings"
)

var compactCode *regexp.Regexp

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
