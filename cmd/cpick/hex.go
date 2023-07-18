package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("hex")

	x.Usage = ""
	x.Summary = "Return a hex value with the \"#\""

	x.Description = `
	The *hex* subcommand is used to return the corresponding hex value
	for a color that is selected when cpick is running.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false)
		if err != nil {
			return err
		}

		fmt.Printf("#%v\n", c.Hex)

		return nil
	}
}
