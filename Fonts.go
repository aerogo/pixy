package main

import (
	"strings"

	"github.com/aerojs/aero"
)

func getFontsCSS() string {
	fontsCSS, err := aero.GetWithHeaders("https://fonts.googleapis.com/css?family=Ubuntu", aero.Headers{
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36",
	})

	if err != nil {
		return ""
	}

	fontsCSS = strings.Replace(fontsCSS, "\r", "", -1)
	fontsCSS = strings.Replace(fontsCSS, "\n", "", -1)
	fontsCSS = strings.Replace(fontsCSS, "  ", " ", -1)
	fontsCSS = strings.Replace(fontsCSS, "{ ", "{", -1)
	fontsCSS = strings.Replace(fontsCSS, ": ", ":", -1)
	fontsCSS = strings.Replace(fontsCSS, "; ", ";", -1)
	fontsCSS = strings.Replace(fontsCSS, ", ", ",", -1)

	return fontsCSS
}
