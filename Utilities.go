package pixy

import (
	"os"
	"os/exec"
	"path"

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
