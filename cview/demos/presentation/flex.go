package main

import (
	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

func demoBox(title string) *cview.Box {
	b := cview.NewBox()
	b.SetBorder(true)
	b.SetTitle(title)
	return b
}

// Flex demonstrates flexbox layout.
func Flex(nextSlide func()) (title string, content cview.Primitive) {
	modalShown := false
	panels := cview.NewPanels()

	textView := cview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle("Flexible width, twice of middle column")
	textView.SetDoneFunc(func(key tcell.Key) {
		if modalShown {
			nextSlide()
			modalShown = false
		} else {
			panels.ShowPanel("modal")
			modalShown = true
		}
	})

	subFlex := cview.NewFlex()
	subFlex.SetDirection(cview.FlexRow)
	subFlex.AddItem(demoBox("Flexible width"), 0, 1, false)
	subFlex.AddItem(demoBox("Fixed height"), 15, 1, false)
	subFlex.AddItem(demoBox("Flexible height"), 0, 1, false)

	flex := cview.NewFlex()
	flex.AddItem(textView, 0, 2, true)
	flex.AddItem(subFlex, 0, 1, false)
	flex.AddItem(demoBox("Fixed width"), 30, 1, false)

	modal := cview.NewModal()
	modal.SetText("Resize the window to see the effect of the flexbox parameters")
	modal.AddButtons([]string{"Ok"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		panels.HidePanel("modal")
	})

	panels.AddPanel("flex", flex, true, true)
	panels.AddPanel("modal", modal, false, false)
	return "Flex", panels
}
