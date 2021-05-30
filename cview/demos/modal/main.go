// Demo code for the Modal primitive.
package main

import (
	"gitlab.com/tslocum/cview"
)

func main() {
	app := cview.NewApplication()
	app.EnableMouse(true)

	modal := cview.NewModal()
	modal.SetText("Do you want to quit the application?")
	modal.AddButtons([]string{"Quit", "Cancel"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			app.Stop()
		}
	})

	app.SetRoot(modal, false)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
