package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ethanbaker/cpick"
)

func main() {
	var sandbox bool
	sandboxUsage := "Run cpick in sandbox mode. No color will be returned\n"
	flag.BoolVar(&sandbox, "sandbox", false, sandboxUsage)
	flag.BoolVar(&sandbox, "s", false, "(shorthand)")

	var testing bool
	testingUsage := "Run cpick in testing mode. No GUI will be shown, only functions tested."
	flag.BoolVar(&testing, "testing", false, testingUsage)
	flag.BoolVar(&testing, "t", false, "(shorthand)")

	flag.Parse()
	args := flag.Args()

	c, err := cpick.Start(sandbox, testing)
	if err != nil {
		log.Fatal(err)
	}

	if !sandbox && !testing {
		args = append(args, "")
		switch strings.ToLower(args[0]) {
		case "rgb":
			fmt.Printf("%v;%v;%v\n", c.RGB.R, c.RGB.G, c.RGB.B)

		case "hsv":
			fmt.Printf("%v;%v;%v\n", c.HSV.H, c.HSV.S, c.HSV.V)

		case "hsl":
			fmt.Printf("%v;%v;%v\n", c.HSL.H, c.HSL.S, c.HSL.L)

		case "cmyk":
			fmt.Printf("%v;%v;%v;%v\n", c.CMYK.C, c.CMYK.M, c.CMYK.Y, c.CMYK.K)

		case "hex":
			fmt.Printf("#%v\n", c.Hex)

		case "decimal":
			fmt.Printf("%v\n", c.Decimal)

		case "ansi":
			fmt.Println(c.Ansi)

		case "escape":
			fmt.Println("\\033[" + strings.Split(string(c.Ansi), "[")[1])

		case "name":
			fmt.Println(c.Name)

		case "json":
			s, err := json.MarshalIndent(c, "", "    ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(s))

		case "bash":
			name := "custom"
			if len(args) >= 3 {
				name = args[1]
			}

			fmt.Printf("readonly -r %v=$'\\033[38;2;%v;%v;%vm'\n", name, c.RGB.R, c.RGB.G, c.RGB.B)

		case "css":
			tag := "color"
			if len(args) >= 3 {
				tag = args[1]
			}

			fmt.Printf("%v: #%v;\n", tag, c.Hex)

		default:
			fmt.Println(c.Ansi)

		}
	} else if testing {
		fmt.Println(true)
	}
}
