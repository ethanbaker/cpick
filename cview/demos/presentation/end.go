package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

// End shows the final slide.
func End(nextSlide func()) (title string, content cview.Primitive) {
	textView := cview.NewTextView()
	textView.SetDoneFunc(func(key tcell.Key) {
		nextSlide()
	})
	url := "https://gitlab.com/tslocum/cview"
	fmt.Fprint(textView, url)
	return "End", Center(len(url), 1, textView)
}
