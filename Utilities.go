package pixy

import (
	"strings"
)

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
