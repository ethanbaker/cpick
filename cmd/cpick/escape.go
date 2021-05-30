package main

import (
	"fmt"
	"strings"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("escape")

	x.Usage = ""
	x.Summary = "Return the ansi escape code"

	x.Description = `
	The *escape* subcommand is used to return the corresponding ansi escape
	code characters for a color that is selected when cpick is running. 
	This will return a line of characters that can be used to represent 
	a color on a terminal.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false, false)
		if err != nil {
			return err
		}

		fmt.Printf("\\033[" + strings.Split(string(c.Ansi), "[")[1])

		return nil
	}
}
