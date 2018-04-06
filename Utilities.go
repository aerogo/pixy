package pixy

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fatih/color"
)

// Runs goimports tool on the specified file.
func goimports(fileOut string) {
	cmd := exec.Command("goimports", "-w", fileOut)
	goimportsErr := cmd.Start()

	if goimportsErr != nil {
		cmd = exec.Command(os.Getenv("GOPATH"), "bin", path.Join("goimports"), "-w", fileOut)
		goimportsErr = cmd.Start()

		if goimportsErr != nil {
			color.Red("Couldn't execute goimports")
			return
		}
	}
}

// extractParameterNames deletes the type information from a comma-separated list of parameters.
func extractParameterNames(definition string) []string {
	definitions := strings.Split(definition, ",")

	for index, definition := range definitions {
		definition := strings.TrimSpace(definition)
		space := strings.Index(definition, " ")

		if space == -1 {
			definitions[index] = definition
			continue
		}

		definitions[index] = definition[:space]
	}

	return definitions
}
