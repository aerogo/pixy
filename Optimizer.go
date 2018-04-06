package pixy

import (
	"bytes"
	"regexp"
	"strings"
)

var compactCode *regexp.Regexp

const (
	writeStringCall = "_b.WriteString("
)

// init
func init() {
	compactCode = regexp.MustCompile("\\n{2,}")
}

// optimize combines multiple WriteString calls to one.
func optimize(code string) (optimizedCode string, inlined string) {
	lines := strings.Split(code, "\n")
	var lastString bytes.Buffer

	// Count the actual code lines
	lineCount := 0

	for index, line := range lines {
		// Find WriteString call
		pos := strings.Index(line, writeStringCall)

		if pos != -1 {
			if line[pos+len(writeStringCall)] == '"' {
				// Delete this line and save it in a buffer "lastString"
				lastString.WriteString(line[pos+len(writeStringCall)+1 : len(line)-2])
				lines[index] = ""
				continue
			}
		}

		if lastString.Len() > 0 {
			lines[index] = "\t" + writeStringCall + "\"" + lastString.String() + "\")\n" + line
			lastString.Reset()
		}

		lineCount++
	}

	compact := compactCode.ReplaceAllString(strings.Join(lines, "\n"), "\n")

	if lineCount == 1 {
		inlined = strings.Replace(compact, writeStringCall, "return ", 1)
		inlined = strings.Replace(inlined, ")\n", "\n", 1)
	}

	return compact, inlined
}
