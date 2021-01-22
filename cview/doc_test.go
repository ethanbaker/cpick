package cview

import (
	"fmt"

	"github.com/gdamore/tcell"
)

// Example of an application with multiple layouts.
func ExampleNewApplication() {
	// Initialize application.
	app := NewApplication()

	// Create shared TextView.
	sharedTextView := NewTextView().SetText("Widgets may be re-used between multiple layouts.").SetTextAlign(AlignCenter)

	// Create main layout using Grid.
	mainTextView := NewTextView().SetText("This is mainLayout.\n\nPress <Tab> to view aboutLayout.").SetTextAlign(AlignCenter)
	mainLayout := NewGrid()
	mainLayout.AddItem(mainTextView, 0, 0, 1, 1, 0, 0, false)
	mainLayout.AddItem(sharedTextView, 1, 0, 1, 1, 0, 0, false)

	// Create about layout using Grid.
	aboutTextView := NewTextView().SetText("cview muti-layout application example\n\nhttps://gitlab.com/tslocum/cview").SetTextAlign(AlignCenter)
	aboutLayout := NewGrid()
	aboutLayout.AddItem(aboutTextView, 0, 0, 1, 1, 0, 0, false)
	aboutLayout.AddItem(sharedTextView, 1, 0, 1, 1, 0, 0, false)

	// Track the current layout.
	currentLayout := 0

	// Set an input capture function that switches between layouts when the tab
	// key is pressed.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if currentLayout == 0 {
				currentLayout = 1

				app.SetRoot(aboutLayout, true)
			} else {
				currentLayout = 0

				app.SetRoot(mainLayout, true)
			}

			// Return nil to stop propagating the event to any remaining
			// handlers.
			return nil
		}

		// Return the event to continue propagating it.
		return event
	})

	// Run the application.
	if err := app.SetRoot(mainLayout, true).Run(); err != nil {
		panic(err)
	}
}

// Example of an application with mouse support.
func ExampleApplication_EnableMouse() {
	// Initialize application.
	app := NewApplication()

	// Enable mouse support.
	app.EnableMouse(true)

	// Enable double clicks.
	app.SetDoubleClickInterval(StandardDoubleClick)

	// Create a textview.
	tv := NewTextView().SetText("Click somewhere!")

	// Set a mouse capture function which prints where the mouse was clicked.
	app.SetMouseCapture(func(event *tcell.EventMouse, action MouseAction) (*tcell.EventMouse, MouseAction) {
		if action == MouseLeftClick || action == MouseLeftDoubleClick {
			actionLabel := "click"
			if action == MouseLeftDoubleClick {
				actionLabel = "double-click"
			}

			x, y := event.Position()

			fmt.Fprintf(tv, "\nYou %sed at %d,%d! Amazing!", actionLabel, x, y)

			// Return nil to stop propagating the event to any remaining handlers.
			return nil, 0
		}

		// Return the event to continue propagating it.
		return event, action
	})

	// Run the application.
	if err := app.SetRoot(tv, true).Run(); err != nil {
		panic(err)
	}
}
