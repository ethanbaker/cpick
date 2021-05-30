package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("bash")

	x.Usage = "<name>"
	x.Summary = "Return a readonly statement with the color constant as an ansi escape code."

	x.Description = `
	The *bash* subcommand is used to return a bash statement
	for a color that is selected when cpick is running. The
	name parameter passed into the command is the name of the
	variable. If no name is specified, the variable name is
	"custom".`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false, false)
		if err != nil {
			return err
		}

		name := "custom"
		if len(args) > 0 {
			name = args[0]
		}

		fmt.Printf("readonly -r %v=$'\\033[38;2;%v;%v;%vm'\n", name, c.RGB.R, c.RGB.G, c.RGB.B)

		return nil
	}
}
