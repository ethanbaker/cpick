package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("name")

	x.Usage = ""
	x.Summary = "Return the name of the color"

	x.Description = `
	The *name* subcommand is used to return the corresponding name
	for a color that is selected when cpick is running. If the 
	selected color has no name, an empty string will be returned.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false)
		if err != nil {
			return err
		}

		fmt.Printf(c.Name)

		return nil
	}
}
