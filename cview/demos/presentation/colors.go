package main

import (
	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

const colorsText = `You can use color tags almost everywhere to partially change the color of a string.

Simply put a color name or hex string in square brackets to change the following characters' color.

H[green]er[white]e i[yellow]s a[darkcyan]n ex[red]amp[white]le.

The [black:red]tags [black:green]look [black:yellow]like [::u]this:

[cyan[]
[blue:yellow:u[]
[#00ff00[]`

// Colors demonstrates how to use colors.
func Colors(nextSlide func()) (title string, content cview.Primitive) {
	tv := cview.NewTextView()
	tv.SetBorder(true)
	tv.SetTitle("A [red]c[yellow]o[green]l[darkcyan]o[blue]r[darkmagenta]f[red]u[yellow]l[white] [black:red]c[:yellow]o[:green]l[:darkcyan]o[:blue]r[:darkmagenta]f[:red]u[:yellow]l[white:] [::bu]title")
	tv.SetDynamicColors(true)
	tv.SetWordWrap(true)
	tv.SetText(colorsText)
	tv.SetDoneFunc(func(key tcell.Key) {
		nextSlide()
	})

	return "Colors", Center(44, 16, tv)
}
