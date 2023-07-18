package main

import (
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("css")

	x.Usage = "<tag>"
	x.Summary = "Return a css line containing a certain tag with the specified color in hexadecimal format"

	x.Description = `
	The *css* subcommand is used to return a css statement
	for a color that is selected when cpick is running. A 
	specific css tag can be specified, which will return that
	tag as the hexadecimal color. If no tag is specified, the
	tag name is "color".`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false)
		if err != nil {
			return err
		}

		tag := "color"
		if len(args) > 0 {
			tag = args[0]
		}

		fmt.Printf("%v: #%v;\n", tag, c.Hex)

		return nil
	}
}
