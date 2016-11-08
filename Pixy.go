package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

const (
	pixyExtension   = ".pixy"
	stylExtension   = ".styl"
	outputName      = "❖"
	outputExtension = ".go"
)

func main() {
	PackageName = "main"

	var output []string

	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		// Pixy
		case pixyExtension:
			fmt.Println(" "+color.GreenString("❀"), path)

			code := CompileFile(path, false)
			output = append(output, code)

		// Stylus
		case stylExtension:
			fmt.Println(" "+color.GreenString("🖌"), path)
			output, err := exec.Command("stylus", "-p", "-c", path).Output()

			if err != nil {
				color.Red("Couldn't execute stylus. Please run 'npm i -g stylus'.")
				return nil
			}

			color.Yellow(string(output))
		}

		return nil
	})

	bundled := strings.Join(output, "\n\n")
	final := getHeader() + bundled

	outputFile := outputName + outputExtension
	writeErr := ioutil.WriteFile(outputFile, []byte(final), 0644)

	if writeErr != nil {
		color.Red("Can't write to " + outputFile)
		return
	}

	cmd := exec.Command("goimports", "-w", outputFile)
	goimportsErr := cmd.Start()

	if goimportsErr != nil {
		workspaceBin := os.Getenv("GOPATH") + string(os.PathSeparator) + "bin" + string(os.PathSeparator)
		cmd = exec.Command(workspaceBin+"goimports", "-w", outputFile)
		goimportsErr = cmd.Start()

		if goimportsErr != nil {
			color.Red("Couldn't execute goimports")
			return
		}
	}

	fmt.Println()
	fmt.Println(" "+color.CyanString("❖"), "Finished.")
}
