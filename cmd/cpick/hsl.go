package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("hsl")

	x.Usage = ""
	x.Summary = "Return hsl values separated by a semi-colon"

	x.Description = `
	The *hsl* subcommand is used to return the HSL values (0-359 for 
	hue, 0-100 for saturation and lightness) for a color that is selected 
	when cpick is running. The HSL values are separated by semi-colons.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false, false)
		if err != nil {
			return err
		}

		fmt.Printf("%v;%v;%v\n", c.HSL.H, c.HSL.S, c.HSL.L)

		return nil
	}
}
