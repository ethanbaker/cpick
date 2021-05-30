package main

import "gitlab.com/tslocum/cview"

// Center returns a new primitive which shows the provided primitive in its
// center, given the provided primitive's size.
func Center(width, height int, p cview.Primitive) cview.Primitive {
	subFlex := cview.NewFlex()
	subFlex.SetDirection(cview.FlexRow)
	subFlex.AddItem(cview.NewBox(), 0, 1, false)
	subFlex.AddItem(p, height, 1, true)
	subFlex.AddItem(cview.NewBox(), 0, 1, false)

	flex := cview.NewFlex()
	flex.AddItem(cview.NewBox(), 0, 1, false)
	flex.AddItem(subFlex, width, 1, true)
	flex.AddItem(cview.NewBox(), 0, 1, false)

	return flex
}
