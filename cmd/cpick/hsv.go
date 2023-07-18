package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("hsv")

	x.Usage = ""
	x.Summary = "Return HSV values separated by a semi-colon"

	x.Description = `
	The *hsv* subcommand is used to return the HSV values (0-359 for 
	hue, 0-100 for saturation and value) for a color that is selected 
	when cpick is running. The HSV values are separated by semi-colons.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false)
		if err != nil {
			return err
		}

		fmt.Printf("%v;%v;%v\n", c.HSV.H, c.HSV.S, c.HSV.V)

		return nil
	}
}
