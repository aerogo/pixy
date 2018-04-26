package pixy

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

// AddImportPaths adds import paths to the specified file.
func AddImportPaths(fileOut ...string) error {
	args := append([]string{"-w"}, fileOut...)
	cmd := exec.Command("goimports", args...)
	goimportsErr := cmd.Start()

	if goimportsErr != nil {
		cmd = exec.Command(path.Join(os.Getenv("GOPATH"), "bin", "goimports"), args...)
		goimportsErr = cmd.Start()

		if goimportsErr != nil {
			return goimportsErr
		}
	}

	return cmd.Wait()
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
