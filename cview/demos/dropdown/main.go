// Demo code for the DropDown primitive.
package main

import "gitlab.com/tslocum/cview"

func main() {
	app := cview.NewApplication()
	app.EnableMouse(true)

	dropdown := cview.NewDropDown()
	dropdown.SetLabel("Select an option (hit Enter): ")
	dropdown.SetOptions(nil,
		cview.NewDropDownOption("First"),
		cview.NewDropDownOption("Second"),
		cview.NewDropDownOption("Third"),
		cview.NewDropDownOption("Fourth"),
		cview.NewDropDownOption("Fifth"))

	app.SetRoot(dropdown, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
