package main

import (
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("cpick", "rgb", "hsv", "hsl", "cmyk", "hex", "decimal", "ansi", "escape", "name", "json", "css", "bash")
	x.Default = "ansi"
	x.Summary = "An extensive color picker for the terminal!"
	x.Version = "1.2.0"
	x.Author = "Ethan Baker <mail@ethanbaker.dev>"
	x.Git = "github.com/ethanbaker/cpick"
	x.Copyright = "(c) Ethan Baker"
	x.License = "Apache-2.0"

	x.Description = `
	Cpick is an interactive color picker in the terminal. You can run Cpick 
	in any true color terminal, and you can see thousands of unique colors, 
	either from preset values or gradients. Each color has its own formats 
	in many different forms, including as RGB, HSV, CMYK, and more.`
}
