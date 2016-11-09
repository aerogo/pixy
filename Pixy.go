package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aerogo/aero"
	"github.com/fatih/color"
)

const (
	pixyExtension   = ".pixy"
	stylExtension   = ".styl"
	outputName      = "$"
	outputExtension = ".go"
)

// StylusCompileResult ...
type StylusCompileResult struct {
	file string
	css  string
}

func main() {
	// Load config file
	app := aero.New()
	app.Load()

	PackageName = "main"

	var output []string
	var css []string

	styleCount := 0
	cssChannel := make(chan *StylusCompileResult, 1024)

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
			go func() {
				style, err := exec.Command("stylus", "-p", "--import", "styles/config.styl", "-c", path).Output()

				if err != nil {
					color.Red("Couldn't execute stylus.")
					color.Red(err.Error())
					cssChannel <- &StylusCompileResult{
						file: path,
						css:  "",
					}
					return
				}

				cssChannel <- &StylusCompileResult{
					file: path,
					css:  string(style),
				}
			}()

			styleCount++
		}

		return nil
	})

	// Fonts
	fontsCSS := getFontsCSS()

	// CSS
	styles := make(map[string]string)

	for i := 0; i < styleCount; i++ {
		result := <-cssChannel
		styles[result.file] = result.css
	}

	// Ordered styles
	for _, styleName := range app.Config.Styles {
		styleName = "styles/" + styleName + ".styl"
		styleContent := styles[styleName]

		if styleContent != "" {
			fmt.Println(" "+color.GreenString("☼"), styleName)
			css = append(css, styleContent)
			styles[styleName] = ""
		}
	}

	// Unordered styles in styles directory
	for styleName, styleContent := range styles {
		if strings.HasPrefix(styleName, "styles/") && styleContent != "" {
			fmt.Println(" "+color.GreenString("☼"), styleName)
			css = append(css, styleContent)
			styles[styleName] = ""
		}
	}

	// Unordered styles
	for styleName, styleContent := range styles {
		if styleContent != "" {
			fmt.Println(" "+color.GreenString("☼"), styleName)
			css = append(css, styleContent)
			styles[styleName] = ""
		}
	}

	bundledCSS := fontsCSS + strings.Join(css, "")
	bundledCSS = strings.Replace(bundledCSS, "\\", "\\\\", -1)
	bundledCSS = strings.Replace(bundledCSS, "\"", "\\\"", -1)

	cssConstant := "const bundledCSS = \"" + bundledCSS + "\"\n"

	bundled := strings.Join(output, "\n\n")
	final := getHeader() + cssConstant + bundled

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
	fmt.Println("Done.")
}
