package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("rgb")

	x.Usage = ""
	x.Summary = "Return RGB values separated by a semi-colon"

	x.Description = `
	The *rgb* subcommand is used to return the corresponding RGB  
	values (between 0 and 255, inclusive) for a color that is selected
	when cpick is running. The RGB values are separated by semi-colons.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false)
		if err != nil {
			return err
		}

		fmt.Printf("%v;%v;%v\n", c.RGB.R, c.RGB.G, c.RGB.B)

		return nil
	}
}
