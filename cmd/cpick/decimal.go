package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("decimal")

	x.Usage = ""
	x.Summary = "Return a deciaml value"

	x.Description = `
	The *decimal* subcommand is used to return the corresponding decimal 
	value	for a color that is selected when cpick is running.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false, false)
		if err != nil {
			return err
		}

		fmt.Printf("%v\n", c.Decimal)

		return nil
	}
}
