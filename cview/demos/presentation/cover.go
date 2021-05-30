package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

const logo = `
 ======= ===  === === ======== ===  ===  ===
===      ===  === === ===      ===  ===  ===
===      ===  === === ======   ===  ===  ===
===       ======  === ===       ===========
 =======    ==    === ========   ==== ====
`

const (
	subtitle   = `Terminal-based user interface toolkit`
	mouse      = `Navigate with your keyboard or mouse.`
	navigation = `Next slide: Ctrl-N   Previous: Ctrl-P   Exit: Ctrl-C`
)

// Cover returns the cover page.
func Cover(nextSlide func()) (title string, content cview.Primitive) {
	// What's the size of the logo?
	lines := strings.Split(logo, "\n")
	logoWidth := 0
	logoHeight := len(lines)
	for _, line := range lines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}
	logoBox := cview.NewTextView()
	logoBox.SetTextColor(tcell.ColorGreen.TrueColor())
	logoBox.SetDoneFunc(func(key tcell.Key) {
		nextSlide()
	})
	fmt.Fprint(logoBox, logo)

	// Create a frame for the subtitle and navigation infos.
	frame := cview.NewFrame(cview.NewBox())
	frame.SetBorders(0, 0, 0, 0, 0, 0)
	frame.AddText(subtitle, true, cview.AlignCenter, tcell.ColorWhite.TrueColor())
	frame.AddText("", true, cview.AlignCenter, tcell.ColorWhite.TrueColor())
	frame.AddText(mouse, true, cview.AlignCenter, tcell.ColorDarkMagenta.TrueColor())
	frame.AddText(navigation, true, cview.AlignCenter, tcell.ColorDarkMagenta.TrueColor())

	// Create a Flex layout that centers the logo and subtitle.
	subFlex := cview.NewFlex()
	subFlex.AddItem(cview.NewBox(), 0, 1, false)
	subFlex.AddItem(logoBox, logoWidth, 1, true)
	subFlex.AddItem(cview.NewBox(), 0, 1, false)

	flex := cview.NewFlex()
	flex.SetDirection(cview.FlexRow)
	flex.AddItem(cview.NewBox(), 0, 7, false)
	flex.AddItem(subFlex, logoHeight, 1, true)
	flex.AddItem(frame, 0, 10, false)

	return "Start", flex
}
