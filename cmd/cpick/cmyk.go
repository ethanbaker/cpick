package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("cmyk")

	x.Usage = ""
	x.Summary = "Return CMYK values separated by a semi-colon"

	x.Description = `
	The *cmyk* subcommand is used to return the CMYK values (between 0
	and 100, inclusive) for a color that is selected when cpick is 
	running. The CMYK values are separated by semi-colons.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false, false)
		if err != nil {
			return err
		}

		fmt.Printf("%v;%v;%v;%v\n", c.CMYK.C, c.CMYK.M, c.CMYK.Y, c.CMYK.K)

		return nil
	}
}
