package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("ansi")

	x.Usage = ""
	x.Summary = "Return the value of an ansi escape code"

	x.Description = `
	The *ansi* subcommand is used to return the corresponding ansi escape
	code for a color that is selected when cpick is running. This will 
	return an actual color on a terminal.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false, false)
		if err != nil {
			return err
		}

		fmt.Println(c.Ansi)

		return nil
	}
}
