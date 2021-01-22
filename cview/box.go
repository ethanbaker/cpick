package cview

import (
	"sync"

	"github.com/gdamore/tcell"
)

// Box is the base Primitive for all widgets. It has a background color and
// optional surrounding elements such as a border and a title. It does not have
// inner text. Widgets embed Box and draw their text over it.
//
// See demos/box for an example.
type Box struct {
	// The position of the rect.
	x, y, width, height int

	// The inner rect reserved for the box's content.
	innerX, innerY, innerWidth, innerHeight int

	// Border padding.
	paddingTop, paddingBottom, paddingLeft, paddingRight int

	// The box's background color.
	backgroundColor tcell.Color

	// Whether or not a border is drawn, reducing the box's space for content by
	// two in width and height.
	border bool

	// The color of the border.
	borderColor tcell.Color

	// The style attributes of the border.
	borderAttributes tcell.AttrMask

	// The title. Only visible if there is a border, too.
	title string

	// The color of the title.
	titleColor tcell.Color

	// The alignment of the title.
	titleAlign int

	// Provides a way to find out if this box has focus. We always go through
	// this interface because it may be overridden by implementing classes.
	focus Focusable

	// Whether or not this box has focus.
	hasFocus bool

	// Whether or not this box shows its focus.
	showFocus bool

	// An optional capture function which receives a key event and returns the
	// event to be forwarded to the primitive's default input handler (nil if
	// nothing should be forwarded).
	inputCapture func(event *tcell.EventKey) *tcell.EventKey

	// An optional function which is called before the box is drawn.
	draw func(screen tcell.Screen, x, y, width, height int) (int, int, int, int)

	// An optional capture function which receives a mouse event and returns the
	// event to be forwarded to the primitive's default mouse event handler (at
	// least one nil if nothing should be forwarded).
	mouseCapture func(action MouseAction, event *tcell.EventMouse) (MouseAction, *tcell.EventMouse)

	l sync.RWMutex
}

// NewBox returns a Box without a border.
func NewBox() *Box {
	b := &Box{
		width:           15,
		height:          10,
		innerX:          -1, // Mark as uninitialized.
		backgroundColor: Styles.PrimitiveBackgroundColor,
		borderColor:     Styles.BorderColor,
		titleColor:      Styles.TitleColor,
		titleAlign:      AlignCenter,
		showFocus:       true,
	}
	b.focus = b
	return b
}

// SetBorderPadding sets the size of the borders around the box content.
func (b *Box) SetBorderPadding(top, bottom, left, right int) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.paddingTop, b.paddingBottom, b.paddingLeft, b.paddingRight = top, bottom, left, right
	return b
}

// GetRect returns the current position of the rectangle, x, y, width, and
// height.
func (b *Box) GetRect() (int, int, int, int) {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.x, b.y, b.width, b.height
}

// GetInnerRect returns the position of the inner rectangle (x, y, width,
// height), without the border and without any padding. Width and height values
// will clamp to 0 and thus never be negative.
func (b *Box) GetInnerRect() (int, int, int, int) {
	b.l.RLock()
	if b.innerX >= 0 {
		defer b.l.RUnlock()
		return b.innerX, b.innerY, b.innerWidth, b.innerHeight
	}
	b.l.RUnlock()

	x, y, width, height := b.GetRect()
	b.l.RLock()
	if b.border {
		x++
		y++
		width -= 2
		height -= 2
	}
	x, y, width, height = x+b.paddingLeft,
		y+b.paddingTop,
		width-b.paddingLeft-b.paddingRight,
		height-b.paddingTop-b.paddingBottom
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	b.l.RUnlock()
	return x, y, width, height
}

// SetRect sets a new position of the primitive. Note that this has no effect
// if this primitive is part of a layout (e.g. Flex, Grid) or if it was added
// like this:
//
//   application.SetRoot(b, true)
func (b *Box) SetRect(x, y, width, height int) {
	b.l.Lock()
	defer b.l.Unlock()

	b.x = x
	b.y = y
	b.width = width
	b.height = height
	b.innerX = -1 // Mark inner rect as uninitialized.
}

// SetDrawFunc sets a callback function which is invoked after the box primitive
// has been drawn. This allows you to add a more individual style to the box
// (and all primitives which extend it).
//
// The function is provided with the box's dimensions (set via SetRect()). It
// must return the box's inner dimensions (x, y, width, height) which will be
// returned by GetInnerRect(), used by descendent primitives to draw their own
// content.
func (b *Box) SetDrawFunc(handler func(screen tcell.Screen, x, y, width, height int) (int, int, int, int)) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.draw = handler
	return b
}

// GetDrawFunc returns the callback function which was installed with
// SetDrawFunc() or nil if no such function has been installed.
func (b *Box) GetDrawFunc() func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.draw
}

// WrapInputHandler wraps an input handler (see InputHandler()) with the
// functionality to capture input (see SetInputCapture()) before passing it
// on to the provided (default) input handler.
//
// This is only meant to be used by subclassing primitives.
func (b *Box) WrapInputHandler(inputHandler func(*tcell.EventKey, func(p Primitive))) func(*tcell.EventKey, func(p Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p Primitive)) {
		if b.inputCapture != nil {
			event = b.inputCapture(event)
		}
		if event != nil && inputHandler != nil {
			inputHandler(event, setFocus)
		}
	}
}

// InputHandler returns nil.
func (b *Box) InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive)) {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.WrapInputHandler(nil)
}

// SetInputCapture installs a function which captures key events before they are
// forwarded to the primitive's default key event handler. This function can
// then choose to forward that key event (or a different one) to the default
// handler by returning it. If nil is returned, the default handler will not
// be called.
//
// Providing a nil handler will remove a previously existing handler.
//
// Note that this function will not have an effect on primitives composed of
// other primitives, such as Form, Flex, or Grid. Key events are only captured
// by the primitives that have focus (e.g. InputField) and only one primitive
// can have focus at a time. Composing primitives such as Form pass the focus on
// to their contained primitives and thus never receive any key events
// themselves. Therefore, they cannot intercept key events.
func (b *Box) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.inputCapture = capture
	return b
}

// GetInputCapture returns the function installed with SetInputCapture() or nil
// if no such function has been installed.
func (b *Box) GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.inputCapture
}

// WrapMouseHandler wraps a mouse event handler (see MouseHandler()) with the
// functionality to capture mouse events (see SetMouseCapture()) before passing
// them on to the provided (default) event handler.
//
// This is only meant to be used by subclassing primitives.
func (b *Box) WrapMouseHandler(mouseHandler func(MouseAction, *tcell.EventMouse, func(p Primitive)) (bool, Primitive)) func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
	return func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
		if b.mouseCapture != nil {
			action, event = b.mouseCapture(action, event)
		}
		if event != nil && mouseHandler != nil {
			consumed, capture = mouseHandler(action, event, setFocus)
		}
		return
	}
}

// MouseHandler returns nil.
func (b *Box) MouseHandler() func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
	return b.WrapMouseHandler(func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
		if action == MouseLeftClick && b.InRect(event.Position()) {
			setFocus(b)
			consumed = true
		}
		return
	})
}

// SetMouseCapture sets a function which captures mouse events (consisting of
// the original tcell mouse event and the semantic mouse action) before they are
// forwarded to the primitive's default mouse event handler. This function can
// then choose to forward that event (or a different one) by returning it or
// returning a nil mouse event, in which case the default handler will not be
// called.
//
// Providing a nil handler will remove a previously existing handler.
func (b *Box) SetMouseCapture(capture func(action MouseAction, event *tcell.EventMouse) (MouseAction, *tcell.EventMouse)) *Box {
	b.mouseCapture = capture
	return b
}

// InRect returns true if the given coordinate is within the bounds of the box's
// rectangle.
func (b *Box) InRect(x, y int) bool {
	rectX, rectY, width, height := b.GetRect()
	return x >= rectX && x < rectX+width && y >= rectY && y < rectY+height
}

// GetMouseCapture returns the function installed with SetMouseCapture() or nil
// if no such function has been installed.
func (b *Box) GetMouseCapture() func(action MouseAction, event *tcell.EventMouse) (MouseAction, *tcell.EventMouse) {
	return b.mouseCapture
}

// SetBackgroundColor sets the box's background color.
func (b *Box) SetBackgroundColor(color tcell.Color) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.backgroundColor = color
	return b
}

// SetBorder sets the flag indicating whether or not the box should have a
// border.
func (b *Box) SetBorder(show bool) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.border = show
	return b
}

// SetBorderColor sets the box's border color.
func (b *Box) SetBorderColor(color tcell.Color) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.borderColor = color
	return b
}

// SetBorderAttributes sets the border's style attributes. You can combine
// different attributes using bitmask operations:
//
//   box.SetBorderAttributes(tcell.AttrUnderline | tcell.AttrBold)
func (b *Box) SetBorderAttributes(attr tcell.AttrMask) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.borderAttributes = attr
	return b
}

// SetTitle sets the box's title.
func (b *Box) SetTitle(title string) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.title = title
	return b
}

// GetTitle returns the box's current title.
func (b *Box) GetTitle() string {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.title
}

// SetTitleColor sets the box's title color.
func (b *Box) SetTitleColor(color tcell.Color) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.titleColor = color
	return b
}

// SetTitleAlign sets the alignment of the title, one of AlignLeft, AlignCenter,
// or AlignRight.
func (b *Box) SetTitleAlign(align int) *Box {
	b.l.Lock()
	defer b.l.Unlock()

	b.titleAlign = align
	return b
}

// Draw draws this primitive onto the screen.
func (b *Box) Draw(screen tcell.Screen) {
	b.l.Lock()

	// Don't draw anything if there is no space.
	if b.width <= 0 || b.height <= 0 {
		b.l.Unlock()
		return
	}

	def := tcell.StyleDefault

	// Fill background.
	background := def.Background(b.backgroundColor)
	for y := b.y; y < b.y+b.height; y++ {
		for x := b.x; x < b.x+b.width; x++ {
			screen.SetContent(x, y, ' ', nil, background)
		}
	}

	// Draw border.
	if b.border && b.width >= 2 && b.height >= 2 {
		border := background.Foreground(b.borderColor) | tcell.Style(b.borderAttributes)
		var vertical, horizontal, topLeft, topRight, bottomLeft, bottomRight rune

		var hasFocus bool
		if b.focus == b {
			hasFocus = b.hasFocus
		} else {
			hasFocus = b.focus.HasFocus()
		}
		if hasFocus && b.showFocus {
			horizontal = Borders.HorizontalFocus
			vertical = Borders.VerticalFocus
			topLeft = Borders.TopLeftFocus
			topRight = Borders.TopRightFocus
			bottomLeft = Borders.BottomLeftFocus
			bottomRight = Borders.BottomRightFocus
		} else {
			horizontal = Borders.Horizontal
			vertical = Borders.Vertical
			topLeft = Borders.TopLeft
			topRight = Borders.TopRight
			bottomLeft = Borders.BottomLeft
			bottomRight = Borders.BottomRight
		}
		for x := b.x + 1; x < b.x+b.width-1; x++ {
			screen.SetContent(x, b.y, horizontal, nil, border)
			screen.SetContent(x, b.y+b.height-1, horizontal, nil, border)
		}
		for y := b.y + 1; y < b.y+b.height-1; y++ {
			screen.SetContent(b.x, y, vertical, nil, border)
			screen.SetContent(b.x+b.width-1, y, vertical, nil, border)
		}
		screen.SetContent(b.x, b.y, topLeft, nil, border)
		screen.SetContent(b.x+b.width-1, b.y, topRight, nil, border)
		screen.SetContent(b.x, b.y+b.height-1, bottomLeft, nil, border)
		screen.SetContent(b.x+b.width-1, b.y+b.height-1, bottomRight, nil, border)

		// Draw title.
		if b.title != "" && b.width >= 4 {
			printed, _ := Print(screen, b.title, b.x+1, b.y, b.width-2, b.titleAlign, b.titleColor)
			if len(b.title)-printed > 0 && printed > 0 {
				_, _, style, _ := screen.GetContent(b.x+b.width-2, b.y)
				fg, _, _ := style.Decompose()
				Print(screen, string(SemigraphicsHorizontalEllipsis), b.x+b.width-2, b.y, 1, AlignLeft, fg)
			}
		}
	}

	// Call custom draw function.
	if b.draw != nil {
		b.l.Unlock()
		newX, newY, newWidth, newHeight := b.draw(screen, b.x, b.y, b.width, b.height)
		b.l.Lock()
		b.innerX, b.innerY, b.innerWidth, b.innerHeight = newX, newY, newWidth, newHeight
	} else {
		// Remember the inner rect.
		b.innerX = -1
		b.l.Unlock()
		newX, newY, newWidth, newHeight := b.GetInnerRect()
		b.l.Lock()
		b.innerX, b.innerY, b.innerWidth, b.innerHeight = newX, newY, newWidth, newHeight
	}

	// Clamp inner rect to screen.
	width, height := screen.Size()
	if b.innerX < 0 {
		b.innerWidth += b.innerX
		b.innerX = 0
	}
	if b.innerX+b.innerWidth >= width {
		b.innerWidth = width - b.innerX
	}
	if b.innerY+b.innerHeight >= height {
		b.innerHeight = height - b.innerY
	}
	if b.innerY < 0 {
		b.innerHeight += b.innerY
		b.innerY = 0
	}
	if b.innerWidth < 0 {
		b.innerWidth = 0
	}
	if b.innerHeight < 0 {
		b.innerHeight = 0
	}

	b.l.Unlock()
}

// Focus is called when this primitive receives focus.
func (b *Box) Focus(delegate func(p Primitive)) {
	b.l.Lock()
	defer b.l.Unlock()

	b.hasFocus = true
}

// Blur is called when this primitive loses focus.
func (b *Box) Blur() {
	b.l.Lock()
	defer b.l.Unlock()

	b.hasFocus = false
}

// HasFocus returns whether or not this primitive has focus.
func (b *Box) HasFocus() bool {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.hasFocus
}

// GetFocusable returns the item's Focusable.
func (b *Box) GetFocusable() Focusable {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.focus
}