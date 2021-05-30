// Demo code for the Button primitive.
package main

import "gitlab.com/tslocum/cview"

func main() {
	app := cview.NewApplication()
	app.EnableMouse(true)

	button := cview.NewButton("Hit Enter to close")
	button.SetBorder(true)
	button.SetRect(0, 0, 22, 3)
	button.SetSelectedFunc(func() {
		app.Stop()
	})

	app.SetRoot(button, false)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
