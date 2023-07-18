package main

import (
	"encoding/json"
	"fmt"

	"github.com/ethanbaker/cpick"
	"github.com/rwxrob/cmdtab"
)

func init() {
	x := cmdtab.New("json")

	x.Usage = ""
	x.Summary = "Return a json object containing all of the color info"

	x.Description = `
	The *json* subcommand is used to return the corresponding json object
	for a color that is selected when cpick is running.`

	x.Method = func(args []string) error {
		c, err := cpick.Start(false)
		if err != nil {
			return err
		}

		s, err := json.MarshalIndent(c, "", "    ")
		if err != nil {
			return err
		}

		fmt.Println(string(s))

		return nil
	}
}
