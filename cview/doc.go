/*
Package cview implements rich widgets for terminal based user interfaces.

See the demos folder and the example application provided with the
Application.NewApplication documentation for usage examples.

Widgets

The following widgets are available:

  Button - Button which is activated when the user selects it.
  Checkbox - Selectable checkbox for boolean values.
  DropDown - Drop-down selection field.
  Flex - A Flexbox based layout manager.
  Form - Form composed of input fields, drop down selections, checkboxes, and
    buttons.
  Grid - A grid based layout manager.  InputField - Single-line text entry field.
  List - A navigable text list with optional keyboard shortcuts.
  Modal - A centered window with a text message and one or more buttons.
  Pages - A page based layout manager.
  ProgressBar - Indicates the progress of an operation.
  Table - A scrollable display of tabular data. Table cells, rows, or columns
    may also be highlighted.
  TextView - A scrollable window that displays multi-colored text. Text may
    also be highlighted.
  TreeView - A scrollable display for hierarchical data. Tree nodes can be
    highlighted, collapsed, expanded, and more.

Base Primitive

Widgets must implement the Primitive interface. All widgets embed the base
primitive, Box, and thus inherit its functions. This isn't necessarily
required, but it makes more sense than reimplementing Box's functionality in
each widget.

Types

This package is a fork of https://github.com/rivo/tview which is based on
https://github.com/gdamore/tcell. It uses types and constants from tcell
(e.g. colors and keyboard values).

Concurrency

All functions may be called concurrently (they are thread-safe). When called
from multiple threads, functions will block until the application or widget
becomes available. Function calls may be queued with Application.QueueUpdate to
avoid blocking.

Unicode Support

This package supports unicode characters including wide characters.

Mouse Support

Mouse support may be enabled by calling Application.EnableMouse before
Application.Run. See the example application provided with the
Application.EnableMouse documentation.

Mouse events are passed to:

- The handler set with SetMouseCapture, which is reserved for use by application
developers to permanently intercept mouse events.

- The ObserveMouseEvent method of every widget under the mouse, bottom to top.

- The MouseHandler method of the topmost widget under the mouse.

Event handlers may return nil to stop propagation.

Colors

Throughout this package, colors are specified using the tcell.Color type.
Functions such as tcell.GetColor(), tcell.NewHexColor(), and tcell.NewRGBColor()
can be used to create colors from W3C color names or RGB values.

Almost all strings which are displayed can contain color tags. Color tags are
W3C color names or six hexadecimal digits following a hash tag, wrapped in
square brackets. Examples:

  This is a [red]warning[white]!
  The sky is [#8080ff]blue[#ffffff].

A color tag changes the color of the characters following that color tag. This
applies to almost everything from box titles, list text, form item labels, to
table cells. In a TextView, this functionality has to be switched on explicitly.
See the TextView documentation for more information.

Color tags may contain not just the foreground (text) color but also the
background color and additional flags. In fact, the full definition of a color
tag is as follows:

  [<foreground>:<background>:<flags>]

Each of the three fields can be left blank and trailing fields can be omitted.
(Empty square brackets "[]", however, are not considered color tags.) Colors
that are not specified will be left unchanged. A field with just a dash ("-")
means "reset to default".

You can specify the following flags (some flags may not be supported by your
terminal):

  l: blink
  b: bold
  d: dim
  r: reverse (switch foreground and background color)
  u: underline

Examples:

  [yellow]Yellow text
  [yellow:red]Yellow text on red background
  [:red]Red background, text color unchanged
  [yellow::u]Yellow text underlined
  [::bl]Bold, blinking text
  [::-]Colors unchanged, flags reset
  [-]Reset foreground color
  [-:-:-]Reset everything
  [:]No effect
  []Not a valid color tag, will print square brackets as they are

In the rare event that you want to display a string such as "[red]" or
"[#00ff1a]" without applying its effect, you need to put an opening square
bracket before the closing square bracket. Note that the text inside the
brackets will be matched less strictly than region or colors tags. I.e. any
character that may be used in color or region tags will be recognized. Examples:

  [red[]      will be output as [red]
  ["123"[]    will be output as ["123"]
  [#6aff00[[] will be output as [#6aff00[]
  [a#"[[[]    will be output as [a#"[[]
  []          will be output as [] (see color tags above)
  [[]         will be output as [[] (not an escaped tag)

You can use the Escape() function to insert brackets automatically where needed.

Styles

When primitives are instantiated, they are initialized with colors taken from
the global Styles variable. You may change this variable to adapt the look and
feel of the primitives to your preferred style.

Hello World

The following is an example application which shows a box titled "Greetings"
containing the text "Hello, world!":

  package main

  import (
    "gitlab.com/tslocum/cview"
  )

  func main() {
    tv := cview.NewTextView()
    tv.SetText("Hello, world!").
       SetBorder(true).
       SetTitle("Greetings")
    if err := cview.NewApplication().SetRoot(tv, true).Run(); err != nil {
      panic(err)
    }
  }

First, we create a TextView with a border and a title. Then we create an
application, set the TextView as its root primitive, and run the event loop.
The application exits when the application's Stop() function is called or when
Ctrl-C is pressed.

If we have a primitive which consumes key presses, we call the application's
SetFocus() function to redirect all key presses to that primitive. Most
primitives then offer ways to install handlers that allow you to react to any
actions performed on them.

Demos

The "demos" subdirectory contains a demo for each widget, as well as a
presentation which gives an overview of the widgets and how they may be used.
*/
package cview
