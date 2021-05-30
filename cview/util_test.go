package cview

import (
	"github.com/gdamore/tcell/v2"
)

// newTestApp returns a new application connected to a simulation screen.
func newTestApp(root Primitive) (*Application, error) {
	// Initialize simulation screen
	sc := tcell.NewSimulationScreen("UTF-8")
	sc.SetSize(80, 24)

	// Initialize application
	app := NewApplication()
	app.SetScreen(sc)
	app.SetRoot(root, true)

	return app, nil
}
