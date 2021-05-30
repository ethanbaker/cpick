NOTE! This is a clone of cview that was locally edited in order to customize cell padding in table cells. The customizable cell padding for table cells is not a perfect feature in this clone; cpick only needs a few specific changes, so some would-be features, like customizable cell padding with borders, do not work properly. 

In addition, the Table's setSelectedFunc was changed so that the space key no longer triggered this function. This was done purely for a cleaner experience with cpick.

# cview - Terminal-based user interface toolkit
[![GoDoc](https://gitlab.com/tslocum/godoc-static/-/raw/master/badge.svg)](https://docs.rocketnine.space/gitlab.com/tslocum/cview)
[![CI status](https://gitlab.com/tslocum/cview/badges/master/pipeline.svg)](https://gitlab.com/tslocum/cview/commits/master)
[![Donate](https://img.shields.io/liberapay/receives/rocketnine.space.svg?logo=liberapay)](https://liberapay.com/rocketnine.space)

This package is a fork of [tview](https://github.com/rivo/tview).
See [FORK.md](https://gitlab.com/tslocum/cview/blob/master/FORK.md) for more information.

## Demo

`ssh cview.rocketnine.space -p 20000`

[![Recording of presentation demo](https://gitlab.com/tslocum/cview/-/raw/master/cview.svg)](https://gitlab.com/tslocum/cview/tree/master/demos/presentation)

## Features

Available widgets:

- __Input forms__ (including __input/password fields__, __drop-down selections__, __checkboxes__, and __buttons__)
- Navigable multi-color __text views__
- Selectable __lists__ with __context menus__
- Modal __dialogs__
- Horizontal and vertical __progress bars__
- __Grid__, __Flexbox__ and __tabbed panel layouts__
- Sophisticated navigable __table views__
- Flexible __tree views__
- Draggable and resizable __windows__
- An __application__ wrapper

Widgets may be customized and extended to suit any application.

[Mouse support](https://docs.rocketnine.space/gitlab.com/tslocum/cview#hdr-Mouse_Support) is available.

## Applications

A list of applications powered by cview is available via [pkg.go.dev](https://pkg.go.dev/gitlab.com/tslocum/cview?tab=importedby).

## Installation

```bash
go get gitlab.com/tslocum/cview
```

## Hello World

This basic example creates a TextView titled "Hello, World!" and displays it in your terminal:

```go
package main

import (
	"gitlab.com/tslocum/cview"
)

func main() {
	app := cview.NewApplication()
	
	box := cview.NewTextView()
	box.SetBorder(true)
	box.SetTitle("Hello, world!")
	box.SetText("Lorem ipsum dolor sit amet")
	
	app.SetRoot(box, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

Examples are available via [godoc](https://docs.rocketnine.space/gitlab.com/tslocum/cview#pkg-examples)
and in the "demos" subdirectory.

For a presentation highlighting the features of this package, compile and run
the program in the "demos/presentation" subdirectory.

## Documentation

Package documentation is available via [godoc](https://docs.rocketnine.space/gitlab.com/tslocum/cview).

An [introduction tutorial](https://rocketnine.space/post/tview-and-you/) is also available.

## Dependencies

This package is based on [github.com/gdamore/tcell](https://github.com/gdamore/tcell)
(and its dependencies) and [github.com/rivo/uniseg](https://github.com/rivo/uniseg).

## Support

[CONTRIBUTING.md](https://gitlab.com/tslocum/cview/blob/master/CONTRIBUTING.md) describes how to share
issues, suggestions and patches (pull requests).
