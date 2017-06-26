package pixy

import (
	"bytes"
	"regexp"
	"strings"
)

var compactCode *regexp.Regexp

const writeStringCall = "_b.WriteString(\""

// init
func init() {
	compactCode = regexp.MustCompile("\\n{2,}")
}

// optimize combines multiple WriteString calls to one.
func optimize(code string) string {
	lines := strings.Split(code, "\n")
	var lastString bytes.Buffer

	// TODO: Optimize single WriteString calls to a simple return

	for index, line := range lines {
		// Find WriteString call
		pos := strings.Index(line, writeStringCall)

		if pos != -1 {
			// Delete this line and save it in a buffer "lastString"
			lastString.WriteString(line[pos+len(writeStringCall) : len(line)-2])
			lines[index] = ""
			continue
		}

		if lastString.Len() > 0 {
			lines[index] = "\t" + writeStringCall + lastString.String() + "\")\n" + line
			lastString.Reset()
		}
	}

	return compactCode.ReplaceAllString(strings.Join(lines, "\n"), "\n")
}
