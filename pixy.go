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
	outputName      = "❖"
	outputExtension = ".go"
)

func main() {
	PackageName = "main"

	var output []string

	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || filepath.Ext(path) != pixyExtension {
			return nil
		}

		// base := filepath.Base(path)
		// outputPath := componentsDirectory + string(os.PathSeparator) + string(base[:len(base)-len(pixyExtension)]) + outputExtension

		// fmt.Println(path, "->", outputPath)
		// CompileFile(path, outputPath)
		fmt.Println(" "+color.GreenString("❀"), path)

		code := CompileFile(path, false)
		output = append(output, code)

		return nil
	})

	bundled := strings.Join(output, "\n\n")
	final := buildHeader(bundled) + bundled

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
